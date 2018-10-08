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
	mpa := NewStubMeterPoints("a", 0, 100*time.Millisecond)
	mpb := NewStubMeterPoints("b", 5, 200*time.Millisecond)
	mpc := NewStubMeterPoints("c", 10, 50*time.Millisecond)
	err := c.Register(mpa, mpb, mpc)
	if err != nil {
		t.Errorf("collector register error: %v", err)
	}
	metrics := c.Retrieve(ctx, time.Second)
	if a, ok := metrics.Get("a.count"); !ok || a != "1" {
		t.Errorf("illegal value a: %q", a)
	}
	metrics = c.Retrieve(ctx, time.Second)
	if b, ok := metrics.Get("b.count"); !ok || b != "7" {
		t.Errorf("illegal value b: %q", b)
	}
	metrics = c.Retrieve(ctx, time.Second)
	if c, ok := metrics.Get("c.count"); !ok || c != "13" {
		t.Errorf("illegal value c: %q", c)
	}
}

// TestCollectorError tests registering meter points with duplicate ID.
func TestCollectorError(t *testing.T) {
	c := collector.New()
	mpa := NewStubMeterPoints("a", 0, 100*time.Millisecond)
	mpb := NewStubMeterPoints("b", 5, 200*time.Millisecond)
	mpc := NewStubMeterPoints("b", 10, 50*time.Millisecond) // <- Double ID.
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
	mpa := NewStubMeterPoints("a", 0, 2*time.Second)
	mpb := NewStubMeterPoints("b", 5, 2*time.Second)
	mpc := NewStubMeterPoints("c", 10, 2*time.Second)
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

// StubMeterPoints simulates meter points.
type StubMeterPoints struct {
	id    string
	count int
	delay time.Duration
}

// NewStubMeterPoints creates a new stub for tests.
func NewStubMeterPoints(id string, start int, delay time.Duration) *StubMeterPoints {
	smp := &StubMeterPoints{
		id:    id,
		count: start,
		delay: delay,
	}
	return smp
}

// ID implements MeterPoints.
func (smp *StubMeterPoints) ID() string {
	return smp.id
}

// Retrieve implements MeterPoints.
func (smp *StubMeterPoints) Retrieve() <-chan collector.Values {
	valuesC := make(chan collector.Values, 1)
	go func() {
		time.Sleep(smp.delay)
		smp.count++
		valuesC <- collector.Values{"count": strconv.Itoa(smp.count)}
	}()
	return valuesC
}

// EOF
