package utils

import "io"

type stopOnly interface {
	Stop()
}

func CloseQuietly(c io.Closer) {
	if c == nil {
		return
	}
	_ = c.Close()
}

func CloseAll(closers ...io.Closer) {
	for _, c := range closers {
		CloseQuietly(c)
	}
}

func StopQuietly(s stopOnly) {
	if s == nil {
		return
	}
	s.Stop()
}
