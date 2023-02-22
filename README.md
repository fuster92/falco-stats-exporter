# Falco Stats File Prometheus Exporter

This is a Prometheus exporter that scrapes metrics from a Falco stats file and exposes them for scraping by Prometheus.
Since Falco doesn't provide the metrics in a format that Prometheus can scrape, this exporter is needed.
## What is Falco?

Falco is a Cloud Native Computing Foundation (CNCF) sandbox project that is used for runtime security and monitoring of container environments. Falco provides a rules engine that can be used to detect and alert on anomalous activity within a container environment.

## What is a Falco stats file?

Falco generates a stats file that provides various metrics related to the events processed by Falco.

## What does this exporter do?

This exporter reads the Falco stats file and exposes the metrics as Prometheus metrics.
The exporter can be run as a standalone binary or as a Docker container.

## Installation
To build and run the exporter, you will need to have Go installed.
The exporter can be built and run with the following commands:

```shell
git clone https://github.com/fuster92/falco-stats-exporter.git
cd falco-stats-exporter
go build -o falco-stats-exporter main.go
./falco-stats-exporter
```

## Usage

To run the exporter, simply execute the binary or run the Docker container. 
By default, the exporter will listen on port 2112 and tail the Falco stats file located at `/tmp/stats.json`.

You can specify a different port or stats file location by passing the following command-line arguments:


## Metrics

The exporter exposes the following metrics:

- `falco_syscall_events_total`: The total number of system call events processed by Falco.
- `falco_syscall_dropped_events_total`: The number of system call events dropped by Falco.
