package heatpump

import (
	"errors"
	"fmt"
	"strconv"
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

// BINARY_HEADER                            TEMP      FAN    P MO          CHSM   P CS
// 1111001000001101000000111111110000000001 1000 0000 0000 0 1 11 00000000 1000 0 1 10
//
// BINARY_HEADER - 40 bits constant header.
// TEMP - 4 bits representing temperature: 0000-1101 (17-30 deg).
// FAN - 4 bits representing fan state: 0000=AUTO, 0100=1, 0110=2, 1000=3, 1010=4, 1100=5.
// P - 1 bit representing power state: 0=ON, 1=OFF.
// MO - 2 bits representing mode: 00=AUTO, 01=COOL, 10=DEHUMIDIFY, 11=HEAT.
// CHSM - 4 bits checksum, calculated by formula: (TEMP + FAN) % 16.
// CS - 2 bits checksum, calculated by formula: MO XOR 01.
//
// ToBinary encodes the state into a binary string for IR transmission.
func (s *State) ToBinary() (string, error) {
	if s.TargetTemperature == nil {
		return "", errors.New("target temperature is nil")
	}
	targetTemperature := *s.TargetTemperature

	if s.FanSpeed == nil {
		return "", errors.New("fan speed is nil")
	}
	fanSpeed := *s.FanSpeed

	if s.Mode == nil {
		return "", errors.New("mode is nil")
	}
	mode := *s.Mode

	temp := targetTemperature - 17

	var fan int
	if fanSpeed > 0 {
		fan = (fanSpeed/20)*2 + 2
	} else {
		fan = 0 // AUTO
	}

	var p int
	if mode == OffMode {
		p = 1 // Note, it's inverted, 0 is ON, 1 is OFF
	} else {
		p = 0
	}

	var mo int
	switch mode {
	case AutoMode:
		mo = 0
	case CoolMode:
		mo = 1
	case HeatMode:
		mo = 3
	case OffMode:
		mo = 3 // For some reason when power is off, the mode is always HEAT
	}

	chsm := (temp + fan) % 16
	cs := mo ^ 1

	return fmt.Sprintf("%s%04b0000%04b0%01b%02b00000000%04b0%01b%02b", BINARY_HEADER, temp, fan, p, mo, chsm, p, cs), nil
}

// NewStateFromBinary decodes the binary string into a State struct.
func NewStateFromBinary(binary string) (*State, error) {
	var s State

	// Validate binary header
	header := binary[:len(BINARY_HEADER)]
	if header != BINARY_HEADER {
		return nil, fmt.Errorf("error invalid binary header, expected: %s, got: %s", BINARY_HEADER, header)
	}

	// Parse binary parts
	temp, err := strconv.ParseInt(binary[40:44], 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing int from binary part TEMP: %v", err)
	}

	fan, err := strconv.ParseInt(binary[48:52], 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing int from binary part FAN: %v", err)
	}

	p, err := strconv.ParseInt(binary[53:54], 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing int from binary part P: %v", err)
	}

	mo, err := strconv.ParseInt(binary[54:56], 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing int from binary part MO: %v", err)
	}

	chsm, err := strconv.ParseInt(binary[64:68], 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing int from binary part CHSM: %v", err)
	}

	cs, err := strconv.ParseInt(binary[70:72], 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing int from binary part CS: %v", err)
	}

	// Validate checksums
	if (temp+fan)%16 != chsm {
		return nil, fmt.Errorf("CHSM checksum mismatch, expected: %d, got: %d", (temp+fan)%16, chsm)
	}

	if mo^1 != cs {
		return nil, fmt.Errorf("MO checksum mismatch, expected: %d, got: %d", mo^1, cs)
	}

	targetTemperature := int(temp) + 17
	s.TargetTemperature = &targetTemperature

	fanSpeed := int(fan) * 20
	s.FanSpeed = &fanSpeed

	var mode Mode
	// If power is on
	if p == 0 {
		switch mo {
		case 0:
			mode = AutoMode
		case 1:
			mode = CoolMode
		case 3:
			mode = HeatMode
		default:
			return nil, fmt.Errorf("mode must be one of: [0, 1, 3], got: %d", mo)
		}
	} else {
		mode = OffMode
	}
	s.Mode = &mode

	return &s, nil
}

type TemperatureReading struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

type Mode string

const (
	OffMode  Mode = "OFF"
	HeatMode Mode = "HEAT"
	CoolMode Mode = "COOL"
	AutoMode Mode = "AUTO"
)

const BINARY_HEADER = "1111001000001101000000111111110000000001"
