// System Monitor Daemon - Collector - Unit Tests
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

// TestCollectorOK tests a correct working collector.
func TestCollectorOK(t *testing.T) {
	c := collector.New()
	mpa := NewStubMeterPoint("a", 0)
	mpb := NewStubMeterPoint("b", 5)
	mpc := NewStubMeterPoint("c", 10)
	err := c.Register(mpa, mpb, mpc)
	if err != nil {
		t.Errorf("collector register error: %v", err)
	}
	metrics := c.Retrieve(time.Second)
	if a, ok := metrics.Get("a"); !ok || a != "1" {
		t.Errorf("illegal value a: %q", a)
	}
	metrics = c.Retrieve(time.Second)
	if b, ok := metrics.Get("b"); !ok || b != "7" {
		t.Errorf("illegal value b: %q", b)
	}
	metrics = c.Retrieve(time.Second)
	if c, ok := metrics.Get("c"); !ok || c != "13" {
		t.Errorf("illegal value c: %q", c)
	}
}

// TestCollectorError tests registering meter points with duplicate ID.
func TestCollectorError(t *testing.T) {
	c := collector.New()
	mpa := NewStubMeterPoint("a", 0)
	mpb := NewStubMeterPoint("b", 5)
	mpc := NewStubMeterPoint("b", 10) // <- Double ID.
	err := c.Register(mpa, mpb, mpc)
	if err == nil {
		t.Errorf("expected registration error")
	}
	if err.Error() != "error: double IDs (b)" {
		t.Errorf("expected different registration error: %v", err)
	}
}

//--------------------
// STUBS
//--------------------

// StubMeterPoint simulates a meter point.
type StubMeterPoint struct {
	id    string
	count int
}

// NewStubMeterPoint creates a new stub for tests.
func NewStubMeterPoint(id string, start int) *StubMeterPoint {
	smp := &StubMeterPoint{
		id:    id,
		count: start,
	}
	return smp
}

// ID implements MeterPoint.
func (smp *StubMeterPoint) ID() string {
	return smp.id
}

// Retrieve implements MeterPoint.
func (smp *StubMeterPoint) Retrieve() <-chan string {
	countC := make(chan string, 1)
	go func() {
		smp.count++
		countC <- strconv.Itoa(smp.count)
	}()
	return countC
}

// EOF
