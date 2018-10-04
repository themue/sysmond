// System Monitor Daemon - Collector - Disk Meter Point - Unit Tests
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
	tests := []struct {
		fullID string
		id     string
		mount  string
		detail collector.DiskDetail
	}{
		{
			fullID: "sys.disk.total.root",
			id:     "root",
			mount:  "/",
			detail: collector.DiskTotal,
		}, {
			fullID: "sys.disk.used.root",
			id:     "root",
			mount:  "/",
			detail: collector.DiskUsed,
		}, {
			fullID: "sys.disk.available.root",
			id:     "root",
			mount:  "/",
			detail: collector.DiskAvailable,
		},
	}
	for i, test := range tests {
		t.Logf("#%d: testing %q", i, test.fullID)
		dmp := collector.NewDiskMeterPoint(test.id, test.mount, test.detail)
		if dmp.ID() != test.fullID {
			t.Errorf("invalid meter point ID: %q", dmp.ID())
		}
		select {
		case value := <-dmp.Retrieve():
			size, err := strconv.Atoi(value)
			if err != nil {
				t.Errorf("meter point value error: %v", err)
			}
			if size <= 0 {
				t.Errorf("invalid meter point value: %d", size)
			}
		case <-time.After(5 * time.Second):
			t.Errorf("meter point timeout")
		}
	}
}

// TestDiskError tests retrieving test data with invalid parameters.
func TestDiskError(t *testing.T) {
	t.Logf("#1: non-existing mount")
	dmp := collector.NewDiskMeterPoint("foo", "/foo/does/not/exist", collector.DiskTotal)
	if dmp.ID() != "sys.disk.total.foo" {
		t.Errorf("invalid meter point ID: %q", dmp.ID())
	}
	select {
	case value := <-dmp.Retrieve():
		if value != "error: cannot retrieve disk space" {
			t.Errorf("invalid return value: %q", value)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter point timeout")
	}
	t.Logf("#2: invalid detail")
	dmp = collector.NewDiskMeterPoint("root", "/", collector.DiskDetail(12345))
	if dmp.ID() != "sys.disk.invalid.root" {
		t.Errorf("invalid meter point ID: %q", dmp.ID())
	}
	select {
	case value := <-dmp.Retrieve():
		if value != "error: invalid detail" {
			t.Errorf("invalid return value: %q", value)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter point timeout")
	}

}

// EOF
