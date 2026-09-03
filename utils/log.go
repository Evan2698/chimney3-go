package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type DefaultWriter struct{}

func (d DefaultWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func configureLogger(output io.Writer) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(output)
}

func defaultLogWriter() io.Writer {
	return &DefaultWriter{}
}

func init() {
	configureLogger(DefaultWriter())
}

func setlogglobal() io.Writer {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	logPath := "log_" + timestamp + ".txt"
	file, err := os.Create(logPath)
	if err != nil {
		fmt.Print("can not create log file", err)
		return defaultLogWriter()
	}
	configureLogger(io.MultiWriter(os.Stdout, file))
	return io.MultiWriter(os.Stdout, file)
}
