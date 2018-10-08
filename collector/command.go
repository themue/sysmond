// System Monitor Daemon - Collector - Command Meter Points
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
	"strings"
)

//--------------------
// COMMAND METER POINTS
//--------------------

// CommandMeterPoint retrieves the single all lines returned by the
// configured command, which typically is a shell script. The lines
// are enumerated.
type CommandMeterPoints struct {
	id      string
	command string
}

// NewCommandMeterPoints creates a new meter point for a passed command.
func NewCommandMeterPoints(id, command string) *CommandMeterPoints {
	return &CommandMeterPoints{
		id:      id,
		command: command,
	}
}

// ID implements MeterPoints.
func (cmp *CommandMeterPoints) ID() string {
	return cmp.id
}

// Retrieve implements MeterPoints.
func (cmp *CommandMeterPoints) Retrieve() <-chan Values {
	valuesC := make(chan Values, 1)
	go func() {
		out, err := exec.Command(cmp.command).Output()
		if err != nil {
			errMsg := fmt.Sprintf("error: cannot execute command: %v", err)
			valuesC <- Values{"1": errMsg}
			return
		}
		lines := strings.Split(string(out), "\n")
		values := make(Values, len(lines))
		for i, line := range lines {
			values[fmt.Sprintf("%d", i+1)] = line
		}
		valuesC <- values
	}()
	return valuesC
}

// EOF
