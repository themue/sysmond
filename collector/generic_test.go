// System Monitor Daemon - Collector - Generic Meter Points - Unit Tests
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
	"errors"
	"testing"
	"time"

	"github.com/themue/sysmond/collector"
)

//--------------------
// TESTS
//--------------------

// TestGenericOK tests a positive generic meter points retrieval.
func TestGenericOK(t *testing.T) {
	gmp := collector.NewGenericMeterPoints("ok", func() (collector.Values, error) {
		return collector.Values{
			"first":  "top",
			"second": "ok",
		}, nil
	})
	if gmp.ID() != "ok" {
		t.Errorf("invalid meter points ID: %q", gmp.ID())
	}
	select {
	case values := <-gmp.Retrieve():
		if values["first"] != "top" || values["second"] != "ok" {
			t.Errorf("invalid meter points values: %q", values)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter points retrieval timeout")
	}
}

// TestGenericError tests a negative generic meter point retrieval.
func TestGenericError(t *testing.T) {
	gmp := collector.NewGenericMeterPoints("nok", func() (collector.Values, error) {
		return nil, errors.New("flop")
	})
	if gmp.ID() != "nok" {
		t.Errorf("invalid meter points ID: %q", gmp.ID())
	}
	select {
	case values := <-gmp.Retrieve():
		if len(values) != 1 || values["all"] != "error: flop" {
			t.Errorf("invalid meter points values: %q", values)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter points retrieval timeout")
	}
}

// EOF
