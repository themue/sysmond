// System Monitor Daemon - Collector - Memory Meter Points
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
	"fmt"
	"os"
	"strings"
)

//--------------------
// MEMORY METER POINT
//--------------------

// MemoryMeterPoints retrieves different memory related metrics by reading
// the /proc/meminfo file.
type MemoryMeterPoints struct {
	prefixes map[string]string
}

// NewMemoryMeterPoints creates new meter points for memory values.
func NewMemoryMeterPoints() *MemoryMeterPoints {
	return &MemoryMeterPoints{
		prefixes: map[string]string{
			"MemTotal:":     "total",
			"MemFree:":      "free",
			"MemAvailable:": "available",
			"Active:":       "active",
			"Inactive:":     "inactive",
			"SwapTotal:":    "swap.total",
			"SwapFree:":     "swap.free",
		},
	}
}

// ID implements MeterPoints.
func (mmp *MemoryMeterPoints) ID() string {
	return "sys.mem"
}

// Retrieve implements MeterPoints.
func (mmp *MemoryMeterPoints) Retrieve() <-chan Values {
	valuesC := make(chan Values, 1)
	go func() {
		file, err := os.Open("/proc/meminfo")
		if err != nil {
			valuesC <- Values{"all": fmt.Sprintf("error: %v", err)}
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		values := make(Values)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			fields := strings.Fields(line)
			id, ok := mmp.prefixes[fields[0]]
			if !ok {
				continue
			}
			values[id] = fields[1]
		}
		valuesC <- values
	}()
	return valuesC
}

// EOF
