package core

import (
	"net"
	"time"
)

func SetConnectTimeout(con net.Conn, tm uint32) error {
	if con == nil {
		return nil
	}

	if tm == 0 {
		return con.SetDeadline(time.Time{})
	}

	deadline := time.Now().Add(time.Duration(tm) * time.Second)
	return con.SetDeadline(deadline)
}
