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

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.SetOutput(&DefaultWriter{})
}

func setlogglobal() io.Writer {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var logpath = "log_" + timestamp + ".txt"
	var file io.Writer
	var err1 error
	file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file", err1)
		file = &DefaultWriter{}
	}
	return io.MultiWriter(os.Stdout, file)
}
