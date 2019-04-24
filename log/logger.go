package log

import (
	"fmt"
	"io"
	"strings"
	"time"
	kitlog "github.com/go-kit/kit/log"
)

type Level int

const (
	Info Level = iota
	Warning
	Debug
	Error
	Fatal
)

func (l Level) String() string {
	switch l {
	case Info: return "I"
	case Warning: return "W"
	case Debug: return "D"
	case Error: return "E"
	case Fatal: return "FATAL"
	default: return "UNKNOWN"
	}
}


type Options struct {
	// Set of levels to log
	Prefix string
	// Setting sync to true wraps the
	Sync bool
	// Defines a filter for each log filter
	Levels map[Level]bool
}

func DefaultOptions() Options {
	return Options{
		Prefix: "",
		Sync: false,
		Levels: map[Level]bool{Info: true, Warning: true, Error: true},
	}
}

type Logger struct {
	w io.Writer
	prefix string
	levels map[Level]bool
}

func NewDefaultLogger(w io.Writer) *Logger {
	return NewLoggerWithOpts(w, DefaultOptions())
}

func NewLoggerWithOpts(w io.Writer, opts Options) *Logger {
	if opts.Sync {
		w = newSyncWriter(w)
	}
	return &Logger{w, opts.Prefix, opts.Levels}
}

// WithPrefix returns a new logger with the prefix appended to the current logger's prefix
func (l Logger) WithPrefix(prefix string) *Logger{
	nextPrefix := strings.Trim(l.prefix + " " + prefix, " ")
	return NewLoggerWithOpts(l.w, Options{
		Prefix: nextPrefix,
		Sync: false,
		Levels: l.levels,
	})
}

func (l *Logger) Info(args...interface{}) {
	l.fprintln(Info, time.Now(), args...)
}

func (l *Logger) Infof(format string, args...interface{}) {
	l.fprintf(Info, time.Now(), format, args...)
}

func (l *Logger) Debug(args...interface{}) {
	l.fprintln(Debug, time.Now(), args...)
}

func (l *Logger) Debugf(format string, args...interface{}) {
	l.fprintf(Debug, time.Now(), format, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.fprintln(Warning, time.Now(), args...)
}

func (l *Logger) Warningf(format string, args...interface{}) {
	l.fprintf(Warning, time.Now(), format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.fprintln(Error, time.Now(), args...)
}

func (l *Logger) Errorf(format string, args...interface{}) {
	l.fprintf(Error, time.Now(), format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	msg := l.completePrefix(Fatal, time.Now()) + " " + fmt.Sprint(args...)
	panic(msg)
}

func (l *Logger) Fatalf(format string, args...interface{}) {
	 msg := l.completePrefix(Fatal, time.Now()) + " " + fmt.Sprintf(format, args...)
	 panic(msg)
}

func (l *Logger) fprintln(level Level, now time.Time, args ...interface{}) {
	if !l.levels[level] && level != Fatal {
		return
	}

	args = append([]interface{}{l.completePrefix(level, now)}, args...)
	_, err := fmt.Fprintln(l.w, args...)
	if err != nil {
		l.Error("LoggingError", err)
	}
}

func (l *Logger) fprintf(level Level, now time.Time, format string, args ...interface{}) {
	l.fprintln(level, now, fmt.Sprintf(format, args...))
}

func (l *Logger) completePrefix(level Level, now time.Time) string {
	return fmt.Sprintf("%s[%s]%s", level, now.Format(time.RFC3339), l.prefix)
}

// newSyncWriter takes an existing io.Writer and wraps a mutex around it to make it safe for use by concurrent goroutines
func newSyncWriter(w io.Writer) io.Writer { return kitlog.NewSyncWriter(w) }
