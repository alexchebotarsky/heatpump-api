package database

import (
	"fmt"

	"github.com/alexchebotarsky/heatpump-api/metrics"
)

const (
	CurrentTemperatureKey = "currentTemperature"
	CurrentHumidityKey    = "currentHumidity"
)

func (d *Database) FetchTemperatureAndHumidity() (temperature float64, humidity float64, err error) {
	temperature, err = d.GetFloat(CurrentTemperatureKey)
	if err != nil {
		return 0, 0, err
	}

	humidity, err = d.GetFloat(CurrentHumidityKey)
	if err != nil {
		return 0, 0, err
	}

	return temperature, humidity, nil
}

func (d *Database) UpdateTemperatureAndHumidity(temperature float64, humidity float64) error {
	err := d.Set(CurrentTemperatureKey, fmt.Sprintf("%.1f", temperature))
	if err != nil {
		return fmt.Errorf("error setting %s in database: %v", CurrentTemperatureKey, err)
	}

	metrics.SetHeatpumpCurrentTemperature(temperature)

	err = d.Set(CurrentHumidityKey, fmt.Sprintf("%.1f", humidity))
	if err != nil {
		return fmt.Errorf("error setting %s in database: %v", CurrentHumidityKey, err)
	}

	metrics.SetHeatpumpCurrentHumidity(humidity)

	return nil
}
