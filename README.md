# System Monitor Daemon

The *System Monitor Daemon* is a small Go based daemon which is intended as
demonstration application for Unix daemons. The application here periodically
collects the values of a number of configured meter points. These values can
be retrieved via HTTP on the configured port. The will be returned in JSON format.

The demo only uses standard library packages and demonstrates different typical
Go idioms for concurrency and synchronisation, for interfaces and higher-order functions.

## Components

### Collector

Major part is the `collector` package defining `MeterPoints` as interface to retrieve
one or more specific meter point values, different implementations of `MeterPoints`
showing the reading of files, the simplified execution of external commands, and
the generic retrieval by user-defined higher-order functions.

The `Collector` retrieves all meter points values in parallel. This process has a
timeout and can also be cancelled by a `context.Context`. The retrieval returns
the `Metrics` which are a set of key/value pairs. The keys are those of the meter 
points ID followed by their individual value IDs, the values the retrieved values.

### Poller

The `Poller` in the `poller` package is a kind of cron for the periodic retrieval
via a collector. The interval can be defined. It is implemented as goroutine using
the actor model to synchronise the access.

### Handler

The `handler` package defines a handler implementing `handler.Handler`. It retrieves
the metrics from a poller and returns these marshalled to JSON after setting the
content-type and a timestamp header.

### SysMonD

Last but not least runs the `sysmond` package the main daemon. It reads a configuration
(so far simulated), creates a handler instance, registers it for the URL path `/metrics`,
and starts the HTTP server in background.

A context timeout or an error are terminating the daemon.

## Open

- Add more tests
- Reading a real configuration file (currently simulated)
- Make the server more flexible (HTTPS, authentication, clean termination)
- More meter points
