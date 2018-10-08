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
	memMP := collector.NewMemoryMeterPoints()
	rootDiskMP := collector.NewDiskMeterPoints("root", "/")
	versionMP := collector.NewGenericMeterPoints("version", func() (collector.Values, error) {
		return collector.Values{"sysmond": version}, nil
	})
	c.Register(memMP, rootDiskMP, versionMP)
	// Return simulated configuration.
	return &Configuration{
		Address:   ":1984",
		Collector: c,
		Interval:  10 * time.Second,
	}, nil
}

// EOF
