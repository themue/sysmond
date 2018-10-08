// System Monitor Daemon - Collector - Disk Meter Points - Unit Tests
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
	"strconv"
	"testing"
	"time"

	"github.com/themue/sysmond/collector"
)

//--------------------
// TESTS
//--------------------

// TestDiskOK tests disk data retrieving with valid parameters.
func TestDiskOK(t *testing.T) {
	testValue := func(values collector.Values, id string) {
		size, err := strconv.Atoi(values[id])
		if err != nil {
			t.Errorf("invalid value: %v", err)
		}
		if size <= 0 {
			t.Errorf("invalid value size: %d", size)
		}
	}
	dmp := collector.NewDiskMeterPoints("root", "/")
	if dmp.ID() != "sys.disk.root" {
		t.Errorf("invalid meter points ID: %q", dmp.ID())
	}
	select {
	case values := <-dmp.Retrieve():
		if len(values) != 3 {
			t.Errorf("invalid number of values: %d", len(values))
		}
		testValue(values, "total")
		testValue(values, "used")
		testValue(values, "available")
	case <-time.After(5 * time.Second):
		t.Errorf("meter points retrieve timeout")
	}
}

// TestDiskError tests retrieving test data with an invalid mount point.
func TestDiskError(t *testing.T) {
	dmp := collector.NewDiskMeterPoints("foo", "/foo/does/not/exist")
	if dmp.ID() != "sys.disk.foo" {
		t.Errorf("invalid meter points ID: %q", dmp.ID())
	}
	select {
	case values := <-dmp.Retrieve():
		if len(values) != 1 || values["all"] != "error: cannot retrieve disk space" {
			t.Errorf("invalid return values: %q", values)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter points retrieve timeout")
	}
}

// EOF
