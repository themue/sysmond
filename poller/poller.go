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

// Poller retrieves system informations via the collector in intervals.
type Poller struct {
	ctx       context.Context
	collector collector.Collector
	interval  time.Duration
}

// New creates a new poller instance.
func New(ctx context.Context, c collector.Collector, i time.Duration) *Poller {
	p := &Poller{
		ctx:       ctx,
		collector: c,
		interval:  i,
	}
	go p.backend()
	return p
}

// backend runs the poller goroutine and calls the collector in intervals.
func (p *Poller) backend() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			metrics := p.collector.Retrieve(p.interval)
		}
	}
}

// EOF
