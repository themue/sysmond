// System Monitor Daemon
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package main

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/themue/sysmond/handler"
)

//--------------------
// CONSTANTS
//--------------------

const version = "v0.1.0"

//--------------------
// RUN
//--------------------

// Run simply configures and runs the server.
func Run(ctx context.Context, cfg *Configuration) <-chan error {
	errC := make(chan error)
	go func() {
		h := handler.New(ctx, cfg.Collector, cfg.Interval)

		http.Handle("/metrics", h)

		errC <- http.ListenAndServe(cfg.Address, nil)
	}()
	return errC
}

//--------------------
// MAIN
//--------------------

// main runs the system monitor daemon.
func main() {
	log.Printf("system monitor daemon %s ...", version)

	// Create a context with timeout to automatically end the demo program.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Read configuration and run the server.
	cfg, err := ReadConfiguration()
	if err != nil {
		log.Fatalf("server configuration error: %v", err)
	}
	errC := Run(ctx, cfg)

	select {
	case <-ctx.Done():
		log.Printf("done!")
	case err = <-errC:
		log.Fatalf("server error: %v", err)
	}
}

// EOF
