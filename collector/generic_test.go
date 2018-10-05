// System Monitor Daemon - Collector - Generic Meter Point - Unit Tests
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

// TestGenericOK tests a positive generic meter point retrieval.
func TestGenericOK(t *testing.T) {
	gmp := collector.NewGenericMeterPoint("ok", func() (string, error) {
		return "top", nil
	})
	if gmp.ID() != "ok" {
		t.Errorf("invalid meter point ID: %q", gmp.ID())
	}
	select {
	case value := <-gmp.Retrieve():
		if value != "top" {
			t.Errorf("invalid meter point value: %q", value)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter point timeout")
	}
}

// TestGenericError tests a negative generic meter point retrieval.
func TestGenericError(t *testing.T) {
	gmp := collector.NewGenericMeterPoint("nok", func() (string, error) {
		return "", errors.New("flop")
	})
	if gmp.ID() != "nok" {
		t.Errorf("invalid meter point ID: %q", gmp.ID())
	}
	select {
	case value := <-gmp.Retrieve():
		if value != "error: flop" {
			t.Errorf("invalid meter point value: %q", value)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("meter point timeout")
	}
}

// EOF
