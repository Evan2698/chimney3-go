package kcpproxy

import (
	"testing"
)

func TestRunKCPRoutineRejectsNilSettings(t *testing.T) {
	if err := RunKCPRoutine(nil, true); err == nil {
		t.Fatal("RunKCPRoutine(nil, true) should reject nil settings")
	}
	if err := RunKCPRoutine(nil, false); err == nil {
		t.Fatal("RunKCPRoutine(nil, false) should reject nil settings")
	}
}
