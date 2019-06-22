package logger

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

// Log provides methods for easy logging. This should be used for each
// running cronjob.
type Log struct {
	jobname string
	start   time.Time

	mu sync.Mutex

	// These two structs both has Write methods to satisfy io.Writer.
	logOut *logOut
	logErr *logErr

	errbuf []byte
	outbuf []byte

	// jsonWriter is used for dbLog
	jsonWriter *JSONBuffer

	dbLog     *log.Logger
	stdoutLog *log.Logger
	stderrLog *log.Logger
}

// LogEntry is one log SQL entry
type LogEntry struct {
	Name     string         `db:"name"`
	ExitCode int            `db:"exit_code"`
	Output   types.JSONText `db:"output"`
	Stamp    time.Time      `db:"stamp"`
	Duration time.Duration  `db:"duration"`
}

// NewLog creates a new logger
func NewLog(job string) *Log {
	var j = NewJSONBuffer()

	log := &Log{
		jobname:    job,
		start:      time.Now(),
		jsonWriter: j,
		dbLog: &log.Logger{
			Handler: json.New(j),
		},
		stdoutLog: &log.Logger{
			Handler: cli.New(os.Stdout),
		},
		stderrLog: &log.Logger{
			Handler: cli.New(os.Stderr),
		},
	}

	log.logOut, log.logErr = newDummyLogs(log)

	return log
}

// Stderr returns a struct that implements io.Writer.
func (l *Log) Stderr() *logErr {
	return l.logErr
}

// Stdout returns a struct that implements io.Writer.
func (l *Log) Stdout() *logOut {
	return l.logOut
}

// Stderr plugs in the io.Writer of stderr
func (l *Log) stderr(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, b := range p {
		if b == '\n' {
			l.flushStderr()
			continue
		}

		l.errbuf = append(l.errbuf, b)
	}

	return len(p), nil
}

func (l *Log) flushStderr() {
	if len(l.errbuf) == 0 {
		return
	}

	// Print the last buffer
	l.stderrLog.WithField("time", time.Now().Format(time.StampMicro)).
		Warn(string(l.errbuf))

	l.dbLog.WithField("output", "stderr").
		Warn(string(l.errbuf))

	// Flush the current buffer
	l.errbuf = l.errbuf[:0]
}

// Stdout plugs in the io.Writer of stderr
func (l *Log) stdout(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, b := range p {
		if b == '\n' {
			l.flushStdout()
			continue
		}

		l.outbuf = append(l.outbuf, b)
	}

	return len(p), nil
}

func (l *Log) flushStdout() {
	if len(l.outbuf) == 0 {
		return
	}

	// Print the last buffer
	l.stderrLog.WithField("time", time.Now().Format(time.StampMicro)).
		Info(string(l.outbuf))

	l.dbLog.WithField("output", "stdout").
		Info(string(l.outbuf))

	// Flush the current buffer
	l.outbuf = l.outbuf[:0]
}

// Finish finalizes the logs and write to the database. The error field should
// be for the command's error.
func (l *Log) Finish(db *factory.DB, err error) (*LogEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Flush the last buffers
	l.flushStderr()
	l.flushStdout()

	var exitCode = 0

	if err != nil {
		// If the error is related to the cron process
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// -2 is the exit code used for errors not related to the cron
			// process
			exitCode = -2
		}

		// Wrap the error
		err = errors.Wrap(err, "Couldn't finish job "+l.jobname)

		// Log the error to the application's stdout
		l.stdoutLog.WithField("output", "stdout").
			Error(err.Error())

		// Log the error to the database as well
		l.dbLog.WithField("output", "stderr").
			Error(err.Error())
	}

	// Create the database entry
	dbLog := &LogEntry{
		Name:     l.jobname,
		ExitCode: exitCode,
		Output:   types.JSONText{},
		Stamp:    l.start,
		Duration: time.Since(l.start),
	}

	// Scan the outputs into the JSON
	if err := dbLog.Output.Scan(l.jsonWriter.String()); err != nil {
		return nil, errors.Wrap(err, "Couldn't scan JSON")
	}

	// Insert the entry into the database
	return dbLog, db.Insert("logs", dbLog)
}
