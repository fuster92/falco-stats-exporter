package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
	"time"
)

type SyscallExporter struct {
	latestEventNumber        safeInt
	latestDroppedEventNumber safeInt
	eventsMetric             *prometheus.Desc
	droppedEventsMetric      *prometheus.Desc
}

func (s *SyscallExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.eventsMetric
	ch <- s.droppedEventsMetric
}

func (s *SyscallExporter) Collect(ch chan<- prometheus.Metric) {
	currentEventNumber := float64(s.latestEventNumber.Get())
	currentDroppedEventNumber := float64(s.latestDroppedEventNumber.Get())
	slog.Debug(
		"collection",
		slog.Group(
			"metrics",
			slog.Float64("currentEventNumber", currentEventNumber),
			slog.Float64(
				"currentDroppedEventNumber",
				currentDroppedEventNumber,
			),
		),
	)
	m1 := prometheus.MustNewConstMetric(
		s.eventsMetric,
		prometheus.CounterValue,
		currentEventNumber,
	)
	m2 := prometheus.MustNewConstMetric(
		s.droppedEventsMetric,
		prometheus.CounterValue,
		currentDroppedEventNumber,
	)
	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
	m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	ch <- m1
	ch <- m2
}

func (s *SyscallExporter) Update(
	latestEventNumber int,
	latestDroppedEventNumber int,
) {
	s.latestEventNumber.Set(latestEventNumber)
	s.latestDroppedEventNumber.Set(latestDroppedEventNumber)
}

func NewFalcoMetricsFilelExporter() *SyscallExporter {
	return &SyscallExporter{
		latestEventNumber:        safeInt{val: 0},
		latestDroppedEventNumber: safeInt{val: 0},
		eventsMetric: prometheus.NewDesc(
			"falco_syscall_events_total",
			"Total number of falco syscall events",
			nil,
			nil,
		),
		droppedEventsMetric: prometheus.NewDesc(
			"falco_syscall_dropped_events_total",
			"Total number of falco syscall events dropped",
			nil,
			nil,
		),
	}
}
