// System Monitor Daemon - Collector - Generic Meter Points
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

// GenericMeterPoints takes a user defined function to retrieve any wanted values.
type GenericMeterPoints struct {
	id       string
	retrieve func() (Values, error)
}

// NewGenericMeterPoints creates new meter points for generic functions.
func NewGenericMeterPoints(id string, r func() (Values, error)) *GenericMeterPoints {
	return &GenericMeterPoints{
		id:       id,
		retrieve: r,
	}
}

// ID implements MeterPoints.
func (gmp *GenericMeterPoints) ID() string {
	return gmp.id
}

// Retrieve implements MeterPoints.
func (gmp *GenericMeterPoints) Retrieve() <-chan Values {
	valuesC := make(chan Values, 1)
	go func() {
		values, err := gmp.retrieve()
		if err != nil {
			valuesC <- Values{"all": fmt.Sprintf("error: %v", err)}
			return
		}
		valuesC <- values
	}()
	return valuesC
}

// EOF
