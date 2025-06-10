package database

import (
	"fmt"

	"github.com/alexchebotarsky/heatpump-api/metrics"
	"github.com/alexchebotarsky/heatpump-api/model/heatpump"
)

const (
	ModeKey              = "mode"
	TargetTemperatureKey = "targetTemperature"
	FanSpeedKey          = "fanSpeed"
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

func snapToNearest(number, snap int) int {
	return ((number + snap/2) / snap) * snap
}
