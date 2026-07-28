package web

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"golang.org/x/net/websocket"

	"github.com/titpetric/platform"

	"github.com/titpetric/phpscript/route"
	"github.com/titpetric/phpscript/runner"
	"github.com/titpetric/phpscript/stdlib"
)

// Module implements platform.Module, serving a PHP-based admin dashboard
// powered by phpscript.
type Module struct {
	platform.UnimplementedModule

	root         fs.FS
	cacheRoot    string
	scriptPath   string
	includeCache *runner.IncludeCache
	exprCache    *runner.ExprCache
}

// NewModule creates the web dashboard module from the given filesystem.
// scriptPath is the directory holding the cron scripts, used by the live "run"
// terminal endpoint.
func NewModule(root fs.FS, scriptPath string) (*Module, error) {
	// Bridge the platform DB env (PLATFORM_DB_CRONTAB) to the phpscript
	// DatabaseDriver env (DB_DSN_CRONTAB) so PHP can open the same database.
	if dsn := os.Getenv("PLATFORM_DB_CRONTAB"); dsn != "" && os.Getenv("DB_DSN_CRONTAB") == "" {
		os.Setenv("DB_DSN_CRONTAB", dsn)
	}

	cacheRoot, err := os.MkdirTemp("", "webcron-frontend-")
	if err != nil {
		return nil, fmt.Errorf("create frontend cache: %w", err)
	}

	return &Module{
		UnimplementedModule: *platform.NewUnimplementedModule("web"),
		root:                &writableFS{source: root, root: cacheRoot},
		cacheRoot:           cacheRoot,
		scriptPath:          scriptPath,
		includeCache:        runner.NewIncludeCache(),
		exprCache:           runner.NewExprCache(),
	}, nil
}

type writableFS struct {
	source fs.FS
	root   string
}

func (f *writableFS) clean(name string) (string, error) {
	if filepath.IsAbs(name) {
		rel, err := filepath.Rel(f.root, name)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", fs.ErrInvalid
		}
		name = rel
	}
	name = path.Clean(filepath.ToSlash(name))
	if !fs.ValidPath(name) {
		return "", fs.ErrInvalid
	}
	return name, nil
}

func (f *writableFS) Open(name string) (fs.File, error) {
	name, err := f.clean(name)
	if err != nil {
		return nil, err
	}
	file, err := os.DirFS(f.root).Open(name)
	if err == nil || !os.IsNotExist(err) {
		return file, err
	}
	return f.source.Open(name)
}

func (f *writableFS) ReadDir(name string) ([]fs.DirEntry, error) {
	name, err := f.clean(name)
	if err != nil {
		return nil, err
	}

	entries := make(map[string]fs.DirEntry)
	for _, filesystem := range []fs.FS{f.source, os.DirFS(f.root)} {
		list, readErr := fs.ReadDir(filesystem, name)
		if readErr != nil && !os.IsNotExist(readErr) {
			return nil, readErr
		}
		for _, entry := range list {
			entries[entry.Name()] = entry
		}
	}

	result := make([]fs.DirEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry)
	}
	slices.SortFunc(result, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return result, nil
}

// Mount registers HTTP routes for the dashboard.
func (m *Module) Mount(ctx context.Context, r platform.Router) error {
	phpMux := http.NewServeMux()
	_, err := route.NewService(m.root, phpMux, route.WithRuntimeFunc(func(rt *runner.Runtime) {
		rt.SetIncludeCache(m.includeCache)
		rt.SetExprCache(m.exprCache)
		stdlib.RegisterFS(rt, m.cacheRoot)
		stdlib.Register(rt)
		registerHelpers(rt)
	}))
	if err != nil {
		return err
	}

	r.Get("/assets/*", m.handleStatic)
	r.Handle("/ws/run/*", websocket.Handler(m.handleRunWS))
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		phpMux.ServeHTTP(w, r)
	}))
	return nil
}

// handleRunWS runs a cron script on demand and streams its combined output to
// the browser over a websocket, so the dashboard can render it in an xterm.js
// terminal. The run is non-interactive: no stdin is forwarded.
func (m *Module) handleRunWS(ws *websocket.Conn) {
	defer ws.Close()

	// Send output as text frames so the browser receives strings directly.
	ws.PayloadType = websocket.TextFrame

	// Resolve the requested job name from the path and guard against traversal.
	name := strings.TrimPrefix(ws.Request().URL.Path, "/ws/run/")
	name = strings.TrimPrefix(path.Clean("/"+name), "/")
	if name == "" || strings.Contains(name, "..") {
		io.WriteString(ws, "invalid job name\r\n")
		return
	}

	script := filepath.Join(m.scriptPath, filepath.FromSlash(name))
	absRoot, err := filepath.Abs(m.scriptPath)
	if err != nil {
		io.WriteString(ws, "server error\r\n")
		return
	}
	absScript, err := filepath.Abs(script)
	if err != nil || !strings.HasPrefix(absScript, absRoot+string(os.PathSeparator)) {
		io.WriteString(ws, "invalid job path\r\n")
		return
	}
	if _, err := os.Stat(script); err != nil {
		fmt.Fprintf(ws, "script not found: %s\r\n", name)
		return
	}

	fmt.Fprintf(ws, "$ ./%s\r\n", name)

	cmd := exec.Command("./" + name)
	cmd.Dir = m.scriptPath
	cmd.Env = append(os.Environ(), "TERM=xterm")

	// A single pipe fed by both stdout and stderr keeps their writes ordered
	// and safe for the copying goroutine (io.PipeWriter is concurrency-safe).
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(ws, "failed to start: %s\r\n", err)
		return
	}

	go func() {
		_ = cmd.Wait()
		_ = pw.Close()
	}()

	_, _ = io.Copy(ws, pr)
	io.WriteString(ws, "\r\n[process finished]\r\n")
}

// handleStatic serves static assets (CSS, JS) from the frontend root.
func (m *Module) handleStatic(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/assets/")
	full := path.Join("assets", path.Clean("/"+relPath))

	data, err := fs.ReadFile(m.root, full)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	switch filepath.Ext(full) {
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	_, _ = w.Write(data)
}

// registerHelpers adds small PHP utility functions that the stdlib doesn't
// provide (date formatting, number formatting, rounding).
func registerHelpers(rt *runner.Runtime) {
	// Override htmlspecialchars to accept any type (the stdlib version
	// requires a string and panics on int64/float64 values from the DB).
	rt.RegisterFunc("htmlspecialchars", func(v any, _ ...any) string {
		s := anyToString(v)
		r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#039;")
		return r.Replace(s)
	})

	rt.RegisterFunc("date", func(format string, ts ...int64) string {
		var t time.Time
		if len(ts) > 0 {
			t = time.Unix(ts[0], 0)
		} else {
			t = time.Now()
		}
		return phpDate(format, t)
	})

	rt.RegisterFunc("time", func() int64 {
		return time.Now().Unix()
	})

	rt.RegisterFunc("strtotime", func(value string) int64 {
		for _, layout := range []string{
			"2006-01-02 15:04:05 -0700 MST",
			"2006-01-02 15:04:05-07:00",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02",
		} {
			if t, err := time.Parse(layout, value); err == nil {
				return t.Unix()
			}
		}
		return 0
	})

	rt.RegisterFunc("number_format", func(v any, args ...int64) string {
		n := toFloat(v)
		dec := 0
		if len(args) > 0 {
			dec = int(args[0])
		}
		return formatNumber(n, dec)
	})

	rt.RegisterFunc("round", func(v any, args ...int64) float64 {
		n := toFloat(v)
		dec := 0
		if len(args) > 0 {
			dec = int(args[0])
		}
		pow := 1.0
		for i := 0; i < dec; i++ {
			pow *= 10
		}
		return float64(int64(n*pow+0.5)) / pow
	})

	rt.RegisterFunc("max", func(a, b any) any {
		if toFloat(a) >= toFloat(b) {
			return a
		}
		return b
	})

	rt.RegisterFunc("min", func(a, b any) any {
		if toFloat(a) <= toFloat(b) {
			return a
		}
		return b
	})
}

// phpDate implements a subset of PHP's date() format: Y-m-d H:i:s and common
// variants used by the dashboard.
func phpDate(format string, t time.Time) string {
	var b strings.Builder
	for i := 0; i < len(format); i++ {
		c := format[i]
		switch c {
		case 'Y':
			b.WriteString(t.Format("2006"))
		case 'm':
			b.WriteString(t.Format("01"))
		case 'd':
			b.WriteString(t.Format("02"))
		case 'H':
			b.WriteString(t.Format("15"))
		case 'i':
			b.WriteString(t.Format("04"))
		case 's':
			b.WriteString(t.Format("05"))
		case 'j':
			b.WriteString(t.Format("2"))
		case 'n':
			b.WriteString(t.Format("1"))
		case 'G':
			b.WriteString(t.Format("3"))
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

func toFloat(v any) float64 {
	switch x := v.(type) {
	case int64:
		return float64(x)
	case int:
		return float64(x)
	case float64:
		return x
	case string:
		var f float64
		_, _ = fmt.Sscanf(x, "%f", &f)
		return f
	default:
		return 0
	}
}

func formatNumber(n float64, dec int) string {
	format := "%." + string(rune('0'+dec)) + "f"
	return fmt.Sprintf(format, n)
}

// anyToString converts any PHP/Go value to its string representation, matching
// PHP's scalar string coercion semantics.
func anyToString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "1"
		}
		return ""
	case int64:
		return fmt.Sprintf("%d", x)
	case int:
		return fmt.Sprintf("%d", x)
	case float64:
		return fmt.Sprintf("%v", x)
	default:
		return fmt.Sprintf("%v", x)
	}
}
