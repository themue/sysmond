// System Monitor Daemon - Configuration
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package main

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/themue/sysmond/collector"
)

//--------------------
// CONFIGURATION
//--------------------

// Configuration contains the configuration to run the system monitor daemon.
type Configuration struct {
	Address   string
	Collector *collector.Collector
	Interval  time.Duration
}

// ReadConfiguration simulates reading a configuration to run the system
// monitor daemon.
func ReadConfiguration() (*Configuration, error) {
	// Configure collector and meter points.
	c := collector.New()
	totalMemMP := collector.NewMemoryMeterPoint(collector.MemTotal)
	freeMemMP := collector.NewMemoryMeterPoint(collector.MemFree)
	availableMemMP := collector.NewMemoryMeterPoint(collector.MemAvailable)
	usedDiskMP := collector.NewDiskMeterPoint("root", "/", collector.DiskUsed)
	availableDiskMP := collector.NewDiskMeterPoint("root", "/", collector.DiskAvailable)
	c.Register(totalMemMP, freeMemMP, availableMemMP, usedDiskMP, availableDiskMP)
	// Return simulated configuration.
	return &Configuration{
		Address:   ":1984",
		Collector: c,
		Interval:  10 * time.Second,
	}, nil
}

// EOF
