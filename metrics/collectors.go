package metrics

import (
	"math"
	"strconv"
	"time"

	"github.com/alexchebotarsky/heatpump-api/model/heatpump"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsHandled = newCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_handled",
		Help: "Handled requests counter and metadata associated with them",
	},
		[]string{"route_name", "status_code"},
	))
	requestsDuration = newCollector(prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "requests_duration",
		Help:    "Time spent processing requests",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	}))

	eventsProcessed = newCollector(prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "events_processed",
		Help: "Handled PubSub events counter and metadata associated with them",
	},
		[]string{"event_name", "status"},
	))
	eventsDuration = newCollector(prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "events_duration",
		Help:    "Time spent processing events",
		Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1.0, 2.5, 5.0, 7.5, 10.0, math.Inf(1)},
	},
		[]string{"event_name"},
	))

	heatpumpMode = newCollector(prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "heatpump_mode",
		Help: "Mode of the heatpump",
	}))
	heatpumpTargetTemperature = newCollector(prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "heatpump_target_temperature",
		Help: "Target temperature of the heatpump",
	}))
	heatpumpFanSpeed = newCollector(prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "heatpump_fan_speed",
		Help: "Fan speed of the heatpump",
	}))

	heatpumpCurrentTemperature = newCollector(prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "heatpump_current_temperature",
		Help: "Current temperature reading of the heatpump",
	}))
	heatpumpCurrentHumidity = newCollector(prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "heatpump_current_humidity",
		Help: "Current humidity reading of the heatpump",
	}))
)

func AddRequestHandled(routeName string, statusCode int) {
	requestsHandled.WithLabelValues(routeName, strconv.Itoa(statusCode)).Inc()
}

func ObserveRequestDuration(duration time.Duration) {
	requestsDuration.Observe(duration.Seconds())
}

func AddEventProcessed(eventName, status string) {
	eventsProcessed.WithLabelValues(eventName, status).Inc()
}

func ObserveEventDuration(eventName string, duration time.Duration) {
	eventsDuration.WithLabelValues(eventName).Observe(duration.Seconds())
}

func SetHeatpumpMode(mode heatpump.Mode) {
	var modeValue float64
	switch mode {
	case heatpump.OffMode:
		modeValue = 0
	case heatpump.HeatMode:
		modeValue = 1
	case heatpump.CoolMode:
		modeValue = 2
	case heatpump.AutoMode:
		modeValue = 3
	default:
		modeValue = -1
	}

	heatpumpMode.Set(modeValue)
}

func SetHeatpumpTargetTemperature(temperature int) {
	heatpumpTargetTemperature.Set(float64(temperature))
}

func SetHeatpumpFanSpeed(fanSpeed int) {
	heatpumpFanSpeed.Set(float64(fanSpeed))
}

func SetHeatpumpCurrentTemperature(temperature float64) {
	heatpumpCurrentTemperature.Set(temperature)
}

func SetHeatpumpCurrentHumidity(humidity float64) {
	heatpumpCurrentHumidity.Set(humidity)
}
