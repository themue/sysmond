# System Monitor Daemon

The *System Monitor Daemon* is a small Go based daemon which is intended as
demonstration application for Unix daemons. The application here periodically
collects the values of a number of configured meter points. These values can
be retrieved via HTTP on the configured port.

## Components

### Collector

### Poller

### Handler

### SysMonD

## Open

- Reading a real configuration file (currently simulated)
- More meter points