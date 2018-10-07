// System Monitor Daemon - Handler
//
// Copyright (C) 2018 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handler

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/themue/sysmond/collector"
	"github.com/themue/sysmond/poller"
)

//--------------------
// METRICS HANDLER
//--------------------

// metricsHandler provides a http.Handler running the metrics server.
type metricsHandler struct {
	collector *collector.Collector
	poller    *poller.Poller
}

// newMetricsHandler returns a new metrics handler instance.
func newMetricsHandler(ctx context.Context, c *collector.Collector, i time.Duration) http.Handler {
	return &metricsHandler{
		collector: c,
		poller:    poller.New(ctx, c, i),
	}
}

// ServeHTTP implements the http.Handler interface.
func (mh *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts, m := mh.poller.Metrics()
	tsb, _ := ts.MarshalText()
	tss := string(tsb)
	b, err := m.Marshal()
	if err != nil {
		errDoc := fmt.Sprintf("{\"error\": \"%v\"}", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Timestamp", tss)
		w.Write([]byte(errDoc))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Timestamp", tss)
	w.Write(b)
}

// EOF
