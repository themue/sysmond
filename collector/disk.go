// System Monitor Daemon - Collector - Disk Meter Point
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
// CONSTANTS
//--------------------

// MemoryDetail defines the detail which is retrieved by the disk meter point.
type DiskDetail int

const (
	DiskTotal DiskDetail = iota + 1
	DiskUsed
	DiskAvailable
)

//--------------------
// DISK METER POINT
//--------------------

// DiskMeterPoint retrieves total, used, and available kilobytes of
// individual file systems.
type DiskMeterPoint struct {
	id     string
	mount  string
	detail DiskDetail
}

// NewDiskMeterPoint creates a new meter point for disk space.
func NewDiskMeterPoint(id, mount string, detail DiskDetail) *DiskMeterPoint {
	dmp := &DiskMeterPoint{
		id:     id,
		mount:  mount,
		detail: detail,
	}
	switch detail {
	case DiskTotal:
		dmp.id = "sys.disk.total." + id
	case DiskUsed:
		dmp.id = "sys.disk.used." + id
	case DiskAvailable:
		dmp.id = "sys.disk.available." + id
	default:
		dmp.id = "sys.disk.invalid." + id
	}
	return dmp
}

// ID implements MeterPoint.
func (dmp *DiskMeterPoint) ID() string {
	return dmp.id
}

// Retrieve implements MeterPoint.
func (dmp *DiskMeterPoint) Retrieve() <-chan string {
	diskC := make(chan string, 1)
	go func() {
		out, err := exec.Command("df", "-Pk", dmp.mount).Output()
		if err != nil {
			diskC <- "error: cannot retrieve disk space"
			return
		}
		lines := strings.Split(string(out), "\n")
		fields := strings.Fields(lines[1])
		switch dmp.detail {
		case DiskTotal:
			diskC <- fields[1]
		case DiskUsed:
			diskC <- fields[2]
		case DiskAvailable:
			diskC <- fields[3]
		default:
			diskC <- "error: invalid detail"
		}
	}()
	return diskC
}

// EOF
