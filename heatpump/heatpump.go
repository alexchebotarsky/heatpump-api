package heatpump

import (
	"fmt"
)

type State struct {
	Mode              *Mode `json:"mode"`
	TargetTemperature *int  `json:"targetTemperature"`
	FanSpeed          *int  `json:"fanSpeed"`
}

func (s *State) Validate() error {
	if s.Mode != nil {
		switch *s.Mode {
		case OffMode, HeatMode, CoolMode, AutoMode:
			// Valid
		default:
			return fmt.Errorf("mode must be one of: [%s, %s, %s, %s], got: %s", OffMode, HeatMode, CoolMode, AutoMode, *s.Mode)
		}
	}

	if s.TargetTemperature != nil {
		if *s.TargetTemperature < 16 || *s.TargetTemperature > 30 {
			return fmt.Errorf("target temperature must be in range [17,30]. got: %d", *s.TargetTemperature)
		}
	}

	if s.FanSpeed != nil {
		if *s.FanSpeed < 0 || *s.FanSpeed > 100 {
			return fmt.Errorf("fan speed must be in range [0,100]. got: %d", *s.FanSpeed)
		}
	}

	return nil
}

type Mode string

const (
	OffMode  Mode = "OFF"
	HeatMode Mode = "HEAT"
	CoolMode Mode = "COOL"
	AutoMode Mode = "AUTO"
)
