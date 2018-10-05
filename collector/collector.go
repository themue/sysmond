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
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

//--------------------
// METER POINT
//--------------------

// MeterPoint defines the interface for individual meter point implementations
// and instances. It is used by the collector to retrieve the according values
// in intervals.
type MeterPoint interface {
	// ID returns the identificator of of the individual meter point.
	ID() string

	// Retrieve returns a channel delivering the polled value. Internal errors
	// have to be returned as string formatted "error: xxx".
	Retrieve() <-chan string
}

//--------------------
// METRICS
//--------------------

// Metrics contains the collected values for marshalling.
type Metrics struct {
	mu     sync.RWMutex
	values map[string]string
}

// NewMetrics creates empty metrics prepared to take the passed number of values.
func NewMetrics(size int) *Metrics {
	return &Metrics{
		values: make(map[string]string, size),
	}
}

// Set sets one value of the metrics.
func (m *Metrics) Set(id, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[id] = value
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
	meterPoints map[string]MeterPoint
}

// New creates a new collector instance.
func New() *Collector {
	return &Collector{
		meterPoints: make(map[string]MeterPoint),
	}
}

// Register adds meter points to the collector. In case of double IDs those
// will be skipped and an error returned.
func (c *Collector) Register(mps ...MeterPoint) error {
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
// at max. the passed duration time, otherwise the value will be
// "error: timeout". All retrievals will be parallel, the wait group
// waits for all retrievals.
func (c *Collector) Retrieve(timeout time.Duration) *Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	metrics := NewMetrics(len(c.meterPoints))
	var wg sync.WaitGroup
	for id, mp := range c.meterPoints {
		wg.Add(1)
		go func(fid string, fmp MeterPoint) {
			defer wg.Done()
			select {
			case value := <-fmp.Retrieve():
				metrics.Set(fid, value)
			case <-time.After(timeout):
				metrics.Set(fid, "error: timeout")
			}
		}(id, mp)
	}
	wg.Wait()
	return metrics
}

// EOF
