// System Monitor Daemon - Collector - CPU Meter Points - Unit Tests
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collector_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/themue/sysmond/collector"
)

//--------------------
// TESTS
//--------------------

// TestCPUOK tests CPU load retrieving with valid parameters.
func TestCPUOK(t *testing.T) {
	testValue := func(values collector.Values, num int, id string) {
		size, err := strconv.ParseFloat(values[fmt.Sprintf("%d.%s", num, id)], 64)
		if err != nil {
			t.Errorf("invalid value: %v", err)
		}
		if size <= 0.0 {
			t.Errorf("invalid value size: %f", size)
		}
	}
	cmp := collector.NewCPUMeterPoints()
	if cmp.ID() != "sys.cpu" {
		t.Errorf("invalid meter points ID: %q", cmp.ID())
	}
	select {
	case values := <-cmp.Retrieve():
		if len(values)%3 != 0 {
			t.Errorf("invalid number of values: %d", len(values))
		}
		cores := len(values) / 3
		for i := 0; i < cores; i++ {
			testValue(values, i, "user")
			testValue(values, i, "system")
			testValue(values, i, "idle")
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter points retrieve timeout")
	}
}

// EOF
