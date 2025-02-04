package middleware

import (
	"context"
	"time"

	"github.com/alexchebotarsky/heatpump-api/metrics"
	"github.com/alexchebotarsky/heatpump-api/processor/event"
)

func Metrics(eventName string, next event.Handler) event.Handler {
	return func(ctx context.Context, payload []byte) error {
		start := time.Now()
		err := next(ctx, payload)
		duration := time.Since(start)

		var status string
		if err != nil {
			status = "ERR"
		} else {
			status = "OK"
		}

		metrics.RecordEventProcessed(eventName, status)
		metrics.ObserveEventDuration(eventName, duration)

		return err
	}
}
