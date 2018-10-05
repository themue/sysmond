// System Monitor Daemon - Collector - Memory Meter Point
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
	"bufio"
	"os"
	"strings"
)

//--------------------
// CONSTANTS
//--------------------

// MemoryDetail defines the detail which is retrieved by the memory meter point.
type MemoryDetail int

const (
	MemTotal MemoryDetail = iota + 1
	MemFree
	MemAvailable
	MemActive
	MemInactive
	MemSwapTotal
	MemSwapFree
)

//--------------------
// MEMORY METER POINT
//--------------------

// MemoryMeterPoint retrieves different memory related metrics by reading
// the /proc/meminfo file.
type MemoryMeterPoint struct {
	detail MemoryDetail
	id     string
	prefix string
}

// NewMemoryMeterPoint creates a new meter point for memory values.
func NewMemoryMeterPoint(detail MemoryDetail) *MemoryMeterPoint {
	mmp := &MemoryMeterPoint{
		detail: detail,
	}
	switch detail {
	case MemTotal:
		mmp.id = "sys.mem.total"
		mmp.prefix = "MemTotal:"
	case MemFree:
		mmp.id = "sys.mem.free"
		mmp.prefix = "MemFree:"
	case MemAvailable:
		mmp.id = "sys.mem.available"
		mmp.prefix = "MemAvailable:"
	case MemActive:
		mmp.id = "sys.mem.active"
		mmp.prefix = "Active:"
	case MemInactive:
		mmp.id = "sys.mem.inactive"
		mmp.prefix = "Inactive:"
	case MemSwapTotal:
		mmp.id = "sys.mem.swap.total"
		mmp.prefix = "SwapTotal:"
	case MemSwapFree:
		mmp.id = "sys.mem.swap.free"
		mmp.prefix = "SwapFree:"
	}
	return mmp
}

// ID implements MeterPoint.
func (mmp *MemoryMeterPoint) ID() string {
	return mmp.id
}

// Retrieve implements MeterPoint.
func (mmp *MemoryMeterPoint) Retrieve() <-chan string {
	memC := make(chan string, 1)
	go func() {
		file, err := os.Open("/proc/meminfo")
		if err != nil {
			memC <- "error: " + err.Error()
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if !strings.HasPrefix(line, mmp.prefix) {
				continue
			}
			fields := strings.Fields(line)
			memC <- fields[1]
			return
		}
		memC <- "error: not found"
	}()
	return memC
}

// EOF
