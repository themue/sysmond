// System Monitor Daemon - Poller - Unit Tests
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package poller_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/themue/sysmond/collector"
	"github.com/themue/sysmond/poller"
)

//--------------------
// TESTS
//--------------------

// TestPollerOK tests a correct working poller.
func TestPollerOK(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	mpsync := collector.NewGenericMeterPoint("sync", func() (string, error) {
		wg.Done()
		return "done", nil
	})
	c := collector.New()
	ctx := context.Background()
	mpa := NewMeterPoint("a")
	mpb := NewMeterPoint("b")
	mpc := NewMeterPoint("c")
	mpd := NewMeterPoint("d")
	err := c.Register(mpa, mpb, mpsync)
	if err != nil {
		t.Errorf("collector register error: %v", err)
	}
	tsBegin := time.Now()
	p := poller.New(ctx, c, 50*time.Millisecond)

	// Wait until three runs are done.
	wg.Wait()

	ts, m := p.Metrics()
	if m == nil {
		t.Errorf("metrics is nil")
		t.FailNow()
	}
	if !ts.After(tsBegin) && !time.Now().After(ts) {
		t.Errorf("illegal metrics timestamp")
	}
	if a, ok := m.Get("a"); !ok || a != "a = 3" {
		t.Errorf("illegal meter point value a: %q", a)
	}
	if b, ok := m.Get("b"); !ok || b != "b = 3" {
		t.Errorf("illegal meter point value b: %q", b)
	}
	if _, ok := m.Get("c"); ok {
		t.Errorf("invalid meter point c")
	}

	// New collector.
	wg.Add(5)
	c = collector.New()
	c.Register(mpb, mpc, mpd, mpsync)
	p.SetCollector(c)

	// Wait until five runs are done.
	wg.Wait()

	ts, m = p.Metrics()
	if m == nil {
		t.Errorf("metrics is nil")
		t.FailNow()
	}
	if !ts.After(tsBegin) && !time.Now().After(ts) {
		t.Errorf("illegal metrics timestamp")
	}
	if _, ok := m.Get("a"); ok {
		t.Errorf("invalid meter point a")
	}
	if b, ok := m.Get("b"); !ok || b != "b = 8" {
		t.Errorf("illegal meter point value b: %q", b)
	}
	if c, ok := m.Get("c"); !ok || c != "c = 5" {
		t.Errorf("illegal meter point value c: %q", c)
	}
	if d, ok := m.Get("d"); !ok || d != "d = 5" {
		t.Errorf("illegal meter point value d: %q", d)
	}
}

//--------------------
// HELPERS
//--------------------

func NewMeterPoint(id string) collector.MeterPoint {
	i := 0
	return collector.NewGenericMeterPoint(id, func() (string, error) {
		i++
		return fmt.Sprintf("%s = %d", id, i), nil
	})
}

// EOF
