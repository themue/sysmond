// System Monitor Daemon - Collector
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collector

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

//--------------------
// VALUE
//--------------------

// Values contains a set of meter point variables and their values.
type Values map[string]string

//--------------------
// METER POINT
//--------------------

// MeterPoints defines the interface for individual meter points implementations
// and instances. It is used by the collector to retrieve the according values
// in intervals.
type MeterPoints interface {
	// ID returns the identificator of of the individual meter points. It
	// defines also the stem for the returned values.
	ID() string

	// Retrieve returns a channel delivering the polled values. Internal errors
	// have to be returned as string value formatted "error: xxx".
	Retrieve() <-chan Values
}

//--------------------
// METRICS
//--------------------

// Metrics contains the collected values for marshalling.
type Metrics struct {
	mu     sync.RWMutex
	values Values
}

// NewMetrics creates empty metrics prepared to take the passed number of values.
func NewMetrics(size int) *Metrics {
	return &Metrics{
		values: make(Values, size),
	}
}

// Set sets one value of the metrics.
func (m *Metrics) Set(id, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[id] = value
}

// Add sets a number of values with a common stem.
func (m *Metrics) Add(stem string, values Values) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, value := range values {
		fullID := fmt.Sprintf("%s.%s", stem, id)
		m.values[fullID] = value
	}
}

// Get reads one value from the metrics.
func (m *Metrics) Get(id string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.values[id]
	return value, ok
}

// Marshal returns the metrics encoded in JSON.
func (m *Metrics) Marshal() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.values)
}

//--------------------
// COLLECTOR
//--------------------

// Collector maintains a number of meter points and retrieves their values
// on demand.
type Collector struct {
	mu          sync.Mutex
	meterPoints map[string]MeterPoints
}

// New creates a new collector instance.
func New() *Collector {
	return &Collector{
		meterPoints: make(map[string]MeterPoints),
	}
}

// Register adds meter points to the collector. In case of double IDs those
// will be skipped and an error returned.
func (c *Collector) Register(mps ...MeterPoints) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var dupes []string
	for _, mp := range mps {
		id := mp.ID()
		if _, ok := c.meterPoints[id]; ok {
			dupes = append(dupes, id)
			continue
		}
		c.meterPoints[id] = mp
	}
	if len(dupes) > 0 {
		return fmt.Errorf("error: double IDs (%s)", strings.Join(dupes, ", "))
	}
	return nil
}

// Retrieve tells the collector to retrieve the metrics. Each has
// at maximum the passed duration time, otherwise the value will be
// "error: timeout". All retrievals will be parallel, the wait group
// waits for all retrievals. The context can cancel the collector
// retrieval as well as all individual goroutines.
func (c *Collector) Retrieve(ctx context.Context, timeout time.Duration) *Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	metrics := NewMetrics(len(c.meterPoints))
	var wg sync.WaitGroup
	for id, mp := range c.meterPoints {
		wg.Add(1)
		go func(fid string, fmp MeterPoints) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				metrics.Set(fid, "error: cancelled")
			case values := <-fmp.Retrieve():
				metrics.Add(fid, values)
			case <-time.After(timeout):
				metrics.Set(fid, "error: timeout")
			}
		}(id, mp)
	}
	wg.Wait()
	return metrics
}

// EOF
