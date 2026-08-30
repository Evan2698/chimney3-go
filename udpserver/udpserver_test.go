package udpserver

import (
	"sync/atomic"
	"testing"
)

func TestStopFlagResetsOnServerStart(t *testing.T) {
	Stop()
	if atomic.LoadInt32(&stop) == 0 {
		t.Fatal("stop flag should be set after Stop()")
	}

	resetStopFlag()
	if atomic.LoadInt32(&stop) != 0 {
		t.Fatal("stop flag should be reset before a new UDP server starts")
	}
}
