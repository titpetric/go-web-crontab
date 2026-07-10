package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/titpetric/platform"

	"github.com/titpetric/phpscript/route"
	"github.com/titpetric/phpscript/runner"
	"github.com/titpetric/phpscript/stdlib"
)

// Module implements platform.Module, serving a PHP-based admin dashboard
// powered by phpscript.
type Module struct {
	platform.UnimplementedModule

	root         string
	includeCache *runner.IncludeCache
	exprCache    *runner.ExprCache
}

// NewModule creates the web dashboard module rooted at the given directory.
func NewModule(root string) *Module {
	// Bridge the platform DB env (PLATFORM_DB_CRONTAB) to the phpscript
	// DatabaseDriver env (DB_DSN_CRONTAB) so PHP can open the same database.
	if dsn := os.Getenv("PLATFORM_DB_CRONTAB"); dsn != "" && os.Getenv("DB_DSN_CRONTAB") == "" {
		os.Setenv("DB_DSN_CRONTAB", dsn)
	}

	return &Module{
		UnimplementedModule: *platform.NewUnimplementedModule("web"),
		root:                root,
		includeCache:        runner.NewIncludeCache(),
		exprCache:           runner.NewExprCache(),
	}
}

// Mount registers HTTP routes for the dashboard.
func (m *Module) Mount(ctx context.Context, r platform.Router) error {
	phpMux := http.NewServeMux()
	_, err := route.NewService(os.DirFS(m.root), phpMux, route.WithRuntimeFunc(func(rt *runner.Runtime) {
		rt.SetIncludeCache(m.includeCache)
		rt.SetExprCache(m.exprCache)
		stdlib.RegisterFS(rt, m.root)
		registerHelpers(rt)
	}))
	if err != nil {
		return err
	}

	r.Get("/assets/*", m.handleStatic)
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		phpMux.ServeHTTP(w, r)
	}))
	return nil
}

// handleStatic serves static assets (CSS, JS) from the frontend root.
func (m *Module) handleStatic(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/assets/")
	full := filepath.Join(m.root, "assets", filepath.Clean("/"+relPath))

	data, err := os.ReadFile(full)
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
