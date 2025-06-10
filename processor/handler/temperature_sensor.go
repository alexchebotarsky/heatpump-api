package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexchebotarsky/heatpump-api/model/heatpump"
	"github.com/alexchebotarsky/heatpump-api/processor/event"
)

type TemperatureAndHumidityUpdater interface {
	UpdateTemperatureAndHumidity(temperature float64, humidity float64) error
}

func TemperatureSensor(updater TemperatureAndHumidityUpdater) event.Handler {
	return func(ctx context.Context, payload []byte) error {
		var reading heatpump.TemperatureReading
		err := json.Unmarshal(payload, &reading)
		if err != nil {
			return fmt.Errorf("error unmarshalling temperature reading: %v", err)
		}

		err = updater.UpdateTemperatureAndHumidity(reading.Temperature, reading.Humidity)
		if err != nil {
			return fmt.Errorf("error updating temperature and humidity: %v", err)
		}

		return nil
	}
}
