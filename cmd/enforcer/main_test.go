package main

import (
    "os"
    "time"

    "testing"
)

func TestRecheckFlag(t *testing.T) {
	os.Args = []string{"enforce", "-recheck=1m"}
	parseFlags() // now we're actually calling the refactored function

	if recheckInterval != time.Minute {
		t.Errorf("Expected recheck interval of %v, got %v", time.Minute, recheckInterval)
	}
}

