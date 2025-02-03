package database

import (
	"fmt"

	"github.com/alexchebotarsky/heatpump-api/heatpump"
)

const (
	ModeKey              = "mode"
	TargetTemperatureKey = "targetTemperature"
	FanSpeedKey          = "fanSpeed"
)

func (d *Database) FetchState() (*heatpump.State, error) {
	var s heatpump.State
	var err error

	modeValue, err := d.GetStr(ModeKey)
	if err != nil {
		return nil, fmt.Errorf("error getting mode from database: %v", err)
	}
	modeEnum := heatpump.Mode(modeValue)
	s.Mode = &modeEnum

	targetTemperature, err := d.GetInt(TargetTemperatureKey)
	if err != nil {
		return nil, fmt.Errorf("error getting target temperature from database: %v", err)
	}
	s.TargetTemperature = &targetTemperature

	fanSpeed, err := d.GetInt(FanSpeedKey)
	if err != nil {
		return nil, fmt.Errorf("error getting fan speed from database: %v", err)
	}
	s.FanSpeed = &fanSpeed

	return &s, nil
}

func (d *Database) UpdateState(state heatpump.State) (*heatpump.State, error) {
	if state.Mode != nil {
		err := d.Set(ModeKey, string(*state.Mode))
		if err != nil {
			return nil, fmt.Errorf("error setting %s in database: %v", ModeKey, err)
		}
	}

	if state.TargetTemperature != nil {
		err := d.Set(TargetTemperatureKey, fmt.Sprintf("%d", *state.TargetTemperature))
		if err != nil {
			return nil, fmt.Errorf("error setting %s in database: %v", TargetTemperatureKey, err)
		}
	}

	if state.FanSpeed != nil {
		err := d.Set(FanSpeedKey, fmt.Sprintf("%d", *state.FanSpeed))
		if err != nil {
			return nil, fmt.Errorf("error setting %s in database: %v", FanSpeedKey, err)
		}
	}

	return d.FetchState()
}
