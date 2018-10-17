// System Monitor Daemon - Collector - CPU Meter Points
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
	"fmt"

	"github.com/shirou/gopsutil/cpu"
)

//--------------------
// CPU METER POINTS
//--------------------

// CPUMeterPoint retrieves the current CPU load.
type CPUMeterPoints struct{}

// NewCPUMeterPoints creates new meter points for CPU load.
func NewCPUMeterPoints() *CPUMeterPoints {
	return &CPUMeterPoints{}
}

// ID implements MeterPoints.
func (cmp *CPUMeterPoints) ID() string {
	return "sys.cpu"
}

// Retrieve implements MeterPoints.
func (cmp *CPUMeterPoints) Retrieve() <-chan Values {
	valuesC := make(chan Values, 1)
	go func() {
		times, err := cpu.Times(true)
		if err != nil {
			valuesC <- Values{"all": "error: cannot retrieve CPU times statistics"}
			return
		}
		values := make(Values, len(times)*3)
		for i, t := range times {
			values[fmt.Sprintf("%d.user", i)] = fmt.Sprintf("%.3f", t.User)
			values[fmt.Sprintf("%d.system", i)] = fmt.Sprintf("%.3f", t.System)
			values[fmt.Sprintf("%d.idle", i)] = fmt.Sprintf("%.3f", t.Idle)
		}
		valuesC <- values
	}()
	return valuesC
}

// EOF
