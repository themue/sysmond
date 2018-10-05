// System Monitor Daemon - Collector - Command Meter Point
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
// COMMAND METER POINT
//--------------------

// CommandMeterPoint retrieves the single first line returned by the
// configured command, which typically is a shell script.
type CommandMeterPoint struct {
	id      string
	command string
}

// NewCommandMeterPoint creates a new meter point for a passed command.
func NewCommandMeterPoint(id, command string) *CommandMeterPoint {
	return &CommandMeterPoint{
		id:      id,
		command: command,
	}
}

// ID implements MeterPoint.
func (cmp *CommandMeterPoint) ID() string {
	return cmp.id
}

// Retrieve implements MeterPoint.
func (cmp *CommandMeterPoint) Retrieve() <-chan string {
	valueC := make(chan string, 1)
	go func() {
		out, err := exec.Command(cmp.command).Output()
		if err != nil {
			valueC <- fmt.Sprintf("error: cannot execute command: %v", err)
			return
		}
		valueC <- string(out)
	}()
	return valueC
}

// EOF
