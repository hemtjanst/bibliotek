package component

import (
	"path"

	"lib.hemtjan.st/v2/class/device"
	"lib.hemtjan.st/v2/class/state"
	"lib.hemtjan.st/v2/platform"
	"lib.hemtjan.st/v2/unit"
)

var _ Settable = (*Sensor)(nil)

type Sensor struct {
	Base
	Template string           `json:"val_tpl,omitempty"`
	Unit     unit.Measurement `json:"unit_of_meas,omitempty"`
	State    state.Class      `json:"stat_cla"`
	StateCh  chan string      `json:"-"`
}

func (s *Sensor) UpdateChannels() []UpdateChannel {
	if s.StateTopic == "" {
		return nil
	}
	return []UpdateChannel{{Topic: s.StateTopic, Channel: s.StateCh}}
}

func NewSensor(name, id string, class device.Class, state state.Class, unit unit.Measurement) *Sensor {
	return &Sensor{
		Base: Base{
			ID:          id,
			Name:        name,
			Platform:    platform.Sensor,
			DeviceClass: class,
			StateTopic:  path.Join("homeassistant", "sensor", id, "state"),
		},
		StateCh: make(chan string),
		Unit:    unit,
		State:   state,
	}
}

func NewTempSensor(name, id string) *Sensor {
	return NewSensor(name, id, device.Temperature, state.Measurement, unit.Celsius)
}

func NewBatterySensor(name, id string) *Sensor {
	return NewSensor(name, id, device.Battery, state.Measurement, unit.Percent)
}
