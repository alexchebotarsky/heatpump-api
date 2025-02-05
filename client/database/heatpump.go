package database

import (
	"fmt"

	"github.com/alexchebotarsky/heatpump-api/metrics"
	"github.com/alexchebotarsky/heatpump-api/model/heatpump"
)

const (
	ModeKey               = "mode"
	TargetTemperatureKey  = "targetTemperature"
	FanSpeedKey           = "fanSpeed"
	CurrentTemperatureKey = "currentTemperature"
	CurrentHumidityKey    = "currentHumidity"
)

func (d *Database) prepareHeatpumpStatements() error {
	mode, err := d.GetStr(ModeKey)
	if err != nil {
		return fmt.Errorf("error getting initial %s from database: %v", ModeKey, err)
	}
	metrics.SetHeatpumpMode(heatpump.Mode(mode))

	targetTemperature, err := d.GetInt(TargetTemperatureKey)
	if err != nil {
		return fmt.Errorf("error getting initial %s from database: %v", TargetTemperatureKey, err)
	}
	metrics.SetHeatpumpTargetTemperature(targetTemperature)

	fanSpeed, err := d.GetInt(FanSpeedKey)
	if err != nil {
		return fmt.Errorf("error getting initial %s from database: %v", FanSpeedKey, err)
	}
	metrics.SetHeatpumpFanSpeed(fanSpeed)

	return nil
}

func (d *Database) FetchHeatpumpState() (*heatpump.State, error) {
	var s heatpump.State
	var err error

	modeValue, err := d.GetStr(ModeKey)
	if err != nil {
		return nil, fmt.Errorf("error getting %s from database: %v", ModeKey, err)
	}
	modeEnum := heatpump.Mode(modeValue)
	s.Mode = &modeEnum

	targetTemperature, err := d.GetInt(TargetTemperatureKey)
	if err != nil {
		return nil, fmt.Errorf("error getting %s from database: %v", TargetTemperatureKey, err)
	}
	s.TargetTemperature = &targetTemperature

	fanSpeed, err := d.GetInt(FanSpeedKey)
	if err != nil {
		return nil, fmt.Errorf("error getting %s from database: %v", FanSpeedKey, err)
	}
	s.FanSpeed = &fanSpeed

	return &s, nil
}

func (d *Database) UpdateHeatpumpState(state *heatpump.State) (*heatpump.State, error) {
	if state.Mode != nil {
		mode := *state.Mode
		err := d.Set(ModeKey, string(mode))
		if err != nil {
			return nil, fmt.Errorf("error setting %s in database: %v", ModeKey, err)
		}
		metrics.SetHeatpumpMode(mode)
	}

	if state.TargetTemperature != nil {
		temperature := *state.TargetTemperature
		err := d.Set(TargetTemperatureKey, fmt.Sprintf("%d", temperature))
		if err != nil {
			return nil, fmt.Errorf("error setting %s in database: %v", TargetTemperatureKey, err)
		}
		metrics.SetHeatpumpTargetTemperature(temperature)
	}

	if state.FanSpeed != nil {
		fanSpeed := snapToNearest(*state.FanSpeed, 20)
		err := d.Set(FanSpeedKey, fmt.Sprintf("%d", fanSpeed))
		if err != nil {
			return nil, fmt.Errorf("error setting %s in database: %v", FanSpeedKey, err)
		}
		metrics.SetHeatpumpFanSpeed(fanSpeed)
	}

	return d.FetchHeatpumpState()
}

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

func snapToNearest(number, snap int) int {
	return ((number + snap/2) / snap) * snap
}
