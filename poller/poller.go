// System Monitor Daemon - Poller
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package poller

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"time"

	"github.com/themue/sysmond/collector"
)

//--------------------
// POLLER
//--------------------

// Poller retrieves system informations via the collector in configurable intervals.
// Stopping is done by the passed context. For serialisation of all accesses the
// backend goroutine works as actor and uses no mutex.
type Poller struct {
	ctx       context.Context
	collector *collector.Collector
	interval  time.Duration
	actionC   chan func()
	timestamp time.Time
	metrics   *collector.Metrics
}

// New creates a new poller instance.
func New(ctx context.Context, c *collector.Collector, i time.Duration) *Poller {
	p := &Poller{
		ctx:       ctx,
		collector: c,
		interval:  i,
		actionC:   make(chan func()),
		timestamp: time.Now(),
	}
	go p.backend()
	return p
}

// SetCollector exchanges the collector. It's not needed, but demonstrates the
// serialised setting of a poller field.
func (p *Poller) SetCollector(c *collector.Collector) {
	p.do(func() {
		p.collector = c
	})
}

// Metrics retrieves the latest metrics and the according timestamp.
func (p *Poller) Metrics() (ts time.Time, m *collector.Metrics) {
	p.do(func() {
		ts = p.timestamp
		m = p.metrics
	})
	return
}

// do lets the actor perform an action in the backend. The wait channel ensures that
// it is performed to avoid race conditions.
func (p *Poller) do(action func()) {
	waitC := make(chan struct{})
	p.actionC <- func() {
		action()
		close(waitC)
	}
	<-waitC
}

// backend runs the poller goroutine and calls the collector in intervals.
func (p *Poller) backend() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case act := <-p.actionC:
			act()
		case <-ticker.C:
			p.timestamp = time.Now()
			p.metrics = p.collector.Retrieve(p.ctx, p.interval)
		}
	}
}

// EOF
