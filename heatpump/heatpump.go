package heatpump

type State struct {
	Mode              *Mode `json:"mode"`
	TargetTemperature *int  `json:"targetTemperature"`
	FanSpeed          *int  `json:"fanSpeed"`
}

type Mode string

const (
	OffMode  Mode = "OFF"
	HeatMode Mode = "HEAT"
	CoolMode Mode = "COOL"
	AutoMode Mode = "AUTO"
)
