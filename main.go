package main

import (
	"flag"
	"fmt"
	"github.com/fuster92/falco-stats-exporter/falco"
	"github.com/nxadm/tail"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func tailMetricFile(metricsFilePath string, errors chan error) (
	<-chan string,
	error,
) {
	t, err := tail.TailFile(
		metricsFilePath, tail.Config{
			Follow: true,
			Logger: tail.DiscardingLogger,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error tailing file: %v", err)
	}
	out := make(chan string)

	go func() {
		for line := range t.Lines {
			out <- line.Text
		}
		err := t.Wait()
		if err != nil {
			slog.Error("error waiting for tail", err)
			errors <- err
		}
	}()
	return out, nil
}

var (
	metricsFile  string
	exporterPort string
	enableDebug  bool
)

func init() {
	flag.StringVar(
		&metricsFile,
		"metricfile",
		"/tmp/metrics.json",
		"falco metricsfile location",
	)
	flag.StringVar(
		&exporterPort,
		"port",
		"2112",
		"port to expose metrics on",
	)
	flag.BoolVar(&enableDebug, "debug", false, "enable debug logging")
	flag.Parse()
}

func setDefaultLogger() {
	var level slog.Level
	if enableDebug {
		level = slog.LevelDebug
	}
	jsonLogger := slog.HandlerOptions{
		Level: level,
	}.NewJSONHandler(os.Stdout)
	logger := slog.New(jsonLogger)

	slog.SetDefault(logger)
}

func main() {
	setDefaultLogger()
	syscallEvents := NewFalcoMetricsFilelExporter()
	errors := make(chan error)
	metricFileLine, err := tailMetricFile(metricsFile, errors)
	if err != nil {
		slog.Error("error tailing metrics file", err)
		os.Exit(1)
	}
	go func() {
		for line := range metricFileLine {
			line = strings.TrimSuffix(line, ",")
			logLine, err := falco.ParseSingleLine(line)
			if err != nil {
				slog.Error("Error parsing line", err)
				continue
			}
			slog.Info("Parsed line", slog.String("line", line))
			syscallEvents.Update(logLine.Cur.Events, logLine.Cur.Drops)
		}
	}()
	prometheus.MustRegister(syscallEvents)
	http.Handle("/metrics", promhttp.Handler())
	slog.Info("starting falco-syscall-exporter", "port", exporterPort)
	go func() {
		errors <- http.ListenAndServe(fmt.Sprintf(":%s", exporterPort), nil)
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	select {
	case err := <-errors:
		slog.Error("fatal error", err)
		os.Exit(1)
	case <-sig:
		slog.Info("Interrupt received, exiting")
	}
}
