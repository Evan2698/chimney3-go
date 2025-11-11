package utils

import (
	"io"
	"log"
)

type DefaultWriter struct{}

func (d DefaultWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Logger is the package-level logger used by utils APIs. Callers can replace
// it with SetLogger or SetLogOutput to redirect logs (helpful for tests).
var Logger *log.Logger

func init() {
	// default logger writes to stdout with a useful timestamp and file info.
	Logger = log.New(&DefaultWriter{}, "", log.Ldate|log.Lmicroseconds|log.Llongfile)
}

// SetLogOutput replaces the output writer for the package logger. If a file
// writer is desired, provide an io.Writer that writes to the file (e.g. os.File).
func SetLogOutput(w io.Writer) {
	if Logger == nil {
		Logger = log.New(w, "", log.Ldate|log.Lmicroseconds|log.Llongfile)
		return
	}
	Logger.SetOutput(w)
}

// SetLogger allows replacing the entire logger (e.g., for tests).
func SetLogger(l *log.Logger) {
	Logger = l
}
