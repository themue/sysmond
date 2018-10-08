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
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/themue/sysmond/collector"
	"github.com/themue/sysmond/poller"
)

//--------------------
// TESTS
//--------------------

// TestPoller tests the working of a poller.
func TestPoller(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	mpsync := collector.NewGenericMeterPoints("sync", func() (collector.Values, error) {
		wg.Done()
		return collector.Values{"wait": "done"}, nil
	})
	testMetric := func(m *collector.Metrics, id, value string) {
		if v, ok := m.Get(id); !ok || v != value {
			t.Errorf("illegal meter point value %q: %q", id, v)
		}
	}
	c := collector.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mpa := NewMeterPoints("a")
	mpb := NewMeterPoints("b")
	mpc := NewMeterPoints("c")
	mpd := NewMeterPoints("d")
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
	testMetric(m, "a.i", "6")
	testMetric(m, "a.j", "20")
	testMetric(m, "b.i", "6")
	testMetric(m, "b.j", "20")
	if _, ok := m.Get("c.i"); ok {
		t.Errorf("invalid meter points value for c.i")
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
	if _, ok := m.Get("a.i"); ok {
		t.Errorf("invalid meter points value for a.i")
	}
	testMetric(m, "b.i", "16")
	testMetric(m, "b.j", "45")
	testMetric(m, "c.i", "10")
	testMetric(m, "c.j", "30")
	testMetric(m, "d.i", "10")
	testMetric(m, "d.j", "30")
}

//--------------------
// HELPERS
//--------------------

func NewMeterPoints(id string) collector.MeterPoints {
	i := 0
	j := 5
	return collector.NewGenericMeterPoints(id, func() (collector.Values, error) {
		i += 2
		j += 5
		return collector.Values{
			"i": strconv.Itoa(i),
			"j": strconv.Itoa(j),
		}, nil
	})
}

// EOF
