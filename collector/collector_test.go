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
	"context"
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
	ctx := context.Background()
	mpa := NewStubMeterPoint("a", 0, 100*time.Millisecond)
	mpb := NewStubMeterPoint("b", 5, 200*time.Millisecond)
	mpc := NewStubMeterPoint("c", 10, 50*time.Millisecond)
	err := c.Register(mpa, mpb, mpc)
	if err != nil {
		t.Errorf("collector register error: %v", err)
	}
	metrics := c.Retrieve(ctx, time.Second)
	if a, ok := metrics.Get("a"); !ok || a != "1" {
		t.Errorf("illegal value a: %q", a)
	}
	metrics = c.Retrieve(ctx, time.Second)
	if b, ok := metrics.Get("b"); !ok || b != "7" {
		t.Errorf("illegal value b: %q", b)
	}
	metrics = c.Retrieve(ctx, time.Second)
	if c, ok := metrics.Get("c"); !ok || c != "13" {
		t.Errorf("illegal value c: %q", c)
	}
}

// TestCollectorError tests registering meter points with duplicate ID.
func TestCollectorError(t *testing.T) {
	c := collector.New()
	mpa := NewStubMeterPoint("a", 0, 100*time.Millisecond)
	mpb := NewStubMeterPoint("b", 5, 200*time.Millisecond)
	mpc := NewStubMeterPoint("b", 10, 50*time.Millisecond) // <- Double ID.
	err := c.Register(mpa, mpb, mpc)
	if err == nil {
		t.Errorf("expected registration error")
	}
	if err.Error() != "error: double IDs (b)" {
		t.Errorf("expected different registration error: %v", err)
	}
}

// TestCollectorCancel tests the cancelling of the retrieval by the context.
func TestCollectorCancel(t *testing.T) {
	c := collector.New()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mpa := NewStubMeterPoint("a", 0, 2*time.Second)
	mpb := NewStubMeterPoint("b", 5, 2*time.Second)
	mpc := NewStubMeterPoint("c", 10, 2*time.Second)
	err := c.Register(mpa, mpb, mpc)
	if err != nil {
		t.Errorf("collector register error: %v", err)
	}
	metrics := c.Retrieve(ctx, 10*time.Second)
	if a, ok := metrics.Get("a"); !ok || a != "error: cancelled" {
		t.Errorf("illegal value a: %q", a)
	}
	if b, ok := metrics.Get("b"); !ok || b != "error: cancelled" {
		t.Errorf("illegal value b: %q", b)
	}
	if c, ok := metrics.Get("c"); !ok || c != "error: cancelled" {
		t.Errorf("illegal value c: %q", c)
	}
}

//--------------------
// STUBS
//--------------------

// StubMeterPoint simulates a meter point.
type StubMeterPoint struct {
	id    string
	count int
	delay time.Duration
}

// NewStubMeterPoint creates a new stub for tests.
func NewStubMeterPoint(id string, start int, delay time.Duration) *StubMeterPoint {
	smp := &StubMeterPoint{
		id:    id,
		count: start,
		delay: delay,
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
		time.Sleep(smp.delay)
		smp.count++
		countC <- strconv.Itoa(smp.count)
	}()
	return countC
}

// EOF
