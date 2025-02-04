package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexchebotarsky/heatpump-api/processor/event"
)

type SensorReading struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type TemperatureAndHumidityUpdater interface {
	UpdateTemperatureAndHumidity(temperature float64, humidity float64) error
}

func DHTSensor(updater TemperatureAndHumidityUpdater) event.Handler {
	return func(ctx context.Context, payload []byte) error {
		var reading SensorReading
		err := json.Unmarshal(payload, &reading)
		if err != nil {
			return fmt.Errorf("error unmarshalling dht reading: %v", err)
		}

		err = updater.UpdateTemperatureAndHumidity(reading.Temperature, reading.Humidity)
		if err != nil {
			return fmt.Errorf("error updating temperature and humidity: %v", err)
		}

		return nil
	}
}
