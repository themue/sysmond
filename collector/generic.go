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
)

//--------------------
// GENERIC METER POINT
//--------------------

// GenericMeterPoint takes a user defined function to retrieve any wanted value.
type GenericMeterPoint struct {
	id       string
	retrieve func() (string, error)
}

// NewGenericMeterPoint creates a new meter point for generic functions.
func NewGenericMeterPoint(id string, r func() (string, error)) *GenericMeterPoint {
	return &GenericMeterPoint{
		id:       id,
		retrieve: r,
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
		value, err := gmp.retrieve()
		if err != nil {
			valueC <- fmt.Sprintf("error: %v", err)
			return
		}
		valueC <- value
	}()
	return valueC
}

// EOF
