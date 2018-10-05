// System Monitor Daemon - Collector - Generic Meter Point
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
	"os/exec"
)

//--------------------
// GENERIC METER POINT
//--------------------

// GenericMeterPoint retrieves the single first line returned by the
// configured command, which typically is a shell script.
type GenericMeterPoint struct {
	id      string
	command string
}

// NewGenericMeterPoint creates a new meter point for disk space.
func NewGenericMeterPoint(id, command string) *GenericMeterPoint {
	return &GenericMeterPoint{
		id:      id,
		command: command,
	}
}

// ID implements MeterPoint.
func (gmp *GenericMeterPoint) ID() string {
	return gmp.id
}

// Retrieve implements MeterPoint.
func (gmp *GenericMeterPoint) Retrieve() <-chan string {
	valueC := make(chan string, 1)
	go func() {
		out, err := exec.Command(gmp.command).Output()
		if err != nil {
			valueC <- fmt.Sprintf("error: cannot execute command: %v", err)
			return
		}
		valueC <- string(out)
	}()
	return valueC
}

// EOF
