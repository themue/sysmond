// System Monitor Daemon - Collector - Marshalling - Unit Tests
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
	"encoding/json"
	"testing"

	"github.com/themue/sysmond/collector"
)

//--------------------
// TESTS
//--------------------

// TestMarshalling tests the marshalling of a metrics.
func TestMarshalling(t *testing.T) {
	m := collector.NewMetrics(5)
	m.Set("a", "1")
	m.Set("b", "two")
	m.Set("c", "drei")
	m.Set("d", "quattre")
	m.Set("e", "cinque")
	b, err := m.Marshal()
	if err != nil {
		t.Errorf("marshalling error: %v", err)
	}
	var tm map[string]string
	err = json.Unmarshal(b, &tm)
	if err != nil {
		t.Errorf("unmarshalling error: %v", err)
	}
	if tm["a"] != "1" {
		t.Errorf("invalid value a: %q", tm["a"])
	}
	if tm["b"] != "two" {
		t.Errorf("invalid value b: %q", tm["b"])
	}
	if tm["c"] != "drei" {
		t.Errorf("invalid value c: %q", tm["c"])
	}
	if tm["d"] != "quattre" {
		t.Errorf("invalid value d: %q", tm["d"])
	}
	if tm["e"] != "cinque" {
		t.Errorf("invalid value e: %q", tm["e"])
	}
}

// EOF
