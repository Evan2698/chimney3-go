package utils

import (
	"time"
)

// Trace records function run information using the package Logger. It
// returns a function that should be deferred by callers to log the elapsed
// time. Example: defer utils.Trace("myFunc")()
func Trace(msg string) func() {
	t1 := time.Now()
	if Logger != nil {
		Logger.Println(" Enter [ " + msg + " ]")
	}
	return func() {
		if Logger != nil {
			Logger.Println(" Leave [ "+msg+" ], takes times:", time.Since(t1))
		}
	}
}
