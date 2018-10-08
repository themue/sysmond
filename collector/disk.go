// System Monitor Daemon - Collector - Disk Meter Points
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
	"os/exec"
	"strings"
)

//--------------------
// DISK METER POINTS
//--------------------

// DiskMeterPoint retrieves total, used, and available kilobytes of
// individual file systems.
type DiskMeterPoints struct {
	id    string
	mount string
}

// NewDiskMeterPoints creates new meter points for disk space.
func NewDiskMeterPoints(id, mount string) *DiskMeterPoints {
	return &DiskMeterPoints{
		id:    "sys.disk." + id,
		mount: mount,
	}
}

// ID implements MeterPoints.
func (dmp *DiskMeterPoints) ID() string {
	return dmp.id
}

// Retrieve implements MeterPoints.
func (dmp *DiskMeterPoints) Retrieve() <-chan Values {
	valuesC := make(chan Values, 1)
	go func() {
		out, err := exec.Command("df", "-Pk", dmp.mount).Output()
		if err != nil {
			valuesC <- Values{"all": "error: cannot retrieve disk space"}
			return
		}
		lines := strings.Split(string(out), "\n")
		fields := strings.Fields(lines[1])
		values := make(Values, 3)
		values["total"] = fields[1]
		values["used"] = fields[2]
		values["available"] = fields[3]
		valuesC <- values
	}()
	return valuesC
}

// EOF
