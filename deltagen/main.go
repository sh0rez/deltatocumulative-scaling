package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
)

const Interval = 15 * time.Second

func main() {
	ctx := context.Background()

	useDelta := func(_ metric.InstrumentKind) metricdata.Temporality {
		return metricdata.DeltaTemporality
	}

	stdout, err := stdoutmetric.New(stdoutmetric.WithTemporalitySelector(useDelta))
	no(err)

	otlp, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithTemporalitySelector(useDelta))
	no(err)

	res, err := resource.New(ctx, resource.WithProcess())
	no(err)

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(stdout, metric.WithInterval(Interval))),
		metric.WithReader(metric.NewPeriodicReader(otlp, metric.WithInterval(Interval))),
		metric.WithResource(res),
	)

	meter := provider.Meter("test.com/deltagen")

	for i := range 100 {
		sum, err := meter.Float64Counter(fmt.Sprintf("rand-%d", i), api.WithDescription("pseudo random garbage"))
		no(err)

		go func() {
			t := time.NewTicker(Interval)
			defer t.Stop()
			for ; true; <-t.C {
				v := float64(rand.Intn(10))
				sum.Add(ctx, v)
			}
		}()
	}

	select {}
}

func no(err error) {
	if err != nil {
		panic(err)
	}
}
