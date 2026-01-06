package component

import (
	"path"

	"lib.hemtjan.st/platform"
	"lib.hemtjan.st/unit"
)

var _ Settable = (*Climate)(nil)

// Climate is an MQTT climate integration
//
// See: https://www.home-assistant.io/integrations/climate.mqtt/
type Climate struct {
	Base
	Modes []Mode `json:"modes"`

	CurrentHumidityTemplate string `json:"current_humidity_template,omitempty"`
	CurrentHumidityTopic    string `json:"current_humidity_topic,omitempty"`

	CurrentTemperatureTemplate string `json:"curr_temp_tpl,omitempty"`
	CurrentTemperatureTopic    string `json:"curr_temp_t,omitempty"`

	EnabledByDefault bool `json:"en"`

	FanModeCommandTemplate string    `json:"fan_mode_cmd_tpl,omitempty"`
	FanModeCommandTopic    string    `json:"fan_mode_cmd_t,omitempty"`
	FanModeStateTemplate   string    `json:"fan_mode_stat_tpl,omitempty"`
	FanModeStateTopic      string    `json:"fan_mode_stat_t,omitempty"`
	FanModes               []FanMode `json:"fan_modes,omitempty"`

	Initial float32 `json:"init,omitzero"`

	JSONAttributesTemplate string `json:"json_attr_tpl,omitempty"`
	JSONAttributesTopic    string `json:"json_attr_t,omitempty"`

	MaxHumidity float32 `json:"max_hum,omitzero"`
	MinHumidity float32 `json:"min_hum,omitzero"`

	MaxTemperature float32 `json:"max_temp,omitzero"`
	MinTemperature float32 `json:"min_temp,omitzero"`

	ModeCommandTemplate string `json:"mode_cmd_tpl,omitempty"`
	ModeCommandTopic    string `json:"mode_cmd_t,omitempty"`
	ModeStateTemplate   string `json:"mode_stat_tpl,omitempty"`
	ModeStateTopic      string `json:"mode_stat_t,omitempty"`

	Optimistic          bool   `json:"opt"`
	PayloadAvailable    string `json:"pl_avail,omitempty"`
	PayloadNotAvailable string `json:"pl_not_avail,omitempty"`
	PayloadOff          string `json:"pl_on,omitempty"`
	PayloadOn           string `json:"pl_off,omitempty"`

	PowerCommandTemplate string `json:"power_command_template,omitempty"`
	PowerCommandTopic    string `json:"power_command_topic,omitempty"`

	Precision float32 `json:"precision,omitzero"`

	PresetModeCommandTemplate string       `json:"pr_mode_cmd_tpl,omitempty"`
	PresetModeCommandTopic    string       `json:"pr_mode_cmd_t,omitempty"`
	PresetModeStateTopic      string       `json:"pr_mode_stat_t,omitempty"`
	PresetModeValueTemplate   string       `json:"pr_mode_val_tpl,omitempty"`
	PresetModes               []PresetMode `json:"pr_modes,omitempty"`

	SwingHorizontalModeCommandTemplateTemplate string                `json:",omitempty"`
	SwingHorizontalModeCommandTopic            string                `json:",omitempty"`
	SwingHorizontalModeStateTemplate           string                `json:",omitempty"`
	SwingHorizontalModeStateTopic              string                `json:",omitempty"`
	SwingHorizontalModes                       []SwingModeHorizontal `json:",omitempty"`
	SwingModeCommandTemplate                   string                `json:",omitempty"`
	SwingModeCommandTopic                      string                `json:",omitempty"`
	SwingModeStateTemplate                     string                `json:",omitempty"`
	SwingModeStateTopic                        string                `json:",omitempty"`
	SwingModes                                 []SwingMode           `json:",omitempty"`

	TargetHumidityCommandTemplate string `json:"hum_cmd_tpl,omitempty"`
	TargetHumidityCommandTopic    string `json:"hum_cmd_t,omitempty"`
	TargetHumidityStateTemplate   string `json:"hum_state_tpl,omitempty"`
	TargetHumidityStateTopic      string `json:"hum_stat_t,omitempty"`

	TemperatureCommandTemplate     string           `json:"temp_cmd_tpl,omitempty"`
	TemperatureCommandTopic        string           `json:"temp_cmd_t,omitempty"`
	TemperatureHighCommandTemplate string           `json:"temp_hi_cmd_tpl,omitempty"`
	TemperatureHighCommandTopic    string           `json:"temp_hi_cmd_t,omitempty"`
	TemperatureHighStateTemplate   string           `json:"temp_hi_stat_tpl,omitempty"`
	TemperatureHighStateTopic      string           `json:"temp_hi_stat_t,omitempty"`
	TemperatureLowCommandTemplate  string           `json:"temp_lo_cmd_tpl,omitempty"`
	TemperatureLowCommandTopic     string           `json:"temp_lo_cmd_t,omitempty"`
	TemperatureLowStateTemplate    string           `json:"temp_lo_stat_tpl,omitempty"`
	TemperatureLowStateTopic       string           `json:"temp_lo_stat_t,omitempty"`
	TemperatureStateTemplate       string           `json:"temp_stat_tpl,omitempty"`
	TemperatureStateTopic          string           `json:"temp_stat_t,omitempty"`
	TemperatureUnit                unit.Measurement `json:"temp_unit,omitempty"`
	TempStep                       float32          `json:"temp_step,omitzero"`

	Template string `json:"val_tpl,omitempty"`

	StateCh chan string `json:"-"`
}

func NewRadiator(id, name string) *Climate {
	prefix := path.Join("homeassistant", "climate", id)

	return &Climate{
		Base: Base{
			ID:         id,
			Name:       name,
			Platform:   platform.Climate,
			BaseTopic:  prefix,
			StateTopic: path.Join(prefix, "state"),
		},
		StateCh: make(chan string),

		EnabledByDefault:        true,
		Modes:                   []Mode{ModeAuto, ModeOff, ModeHeat},
		CurrentTemperatureTopic: path.Join(prefix, "current_temp"),
		MinTemperature:          5.0,
		MaxTemperature:          30.0,
		ModeStateTopic:          path.Join(prefix, "mode", "state"),
		PresetModeStateTopic:    path.Join(prefix, "preset", "state"),
		PresetModes:             []PresetMode{PresetEco, PresetAway, PresetBoost, PresetComfort},
		TemperatureStateTopic:   path.Join(prefix, "temp", "state"),
		TemperatureUnit:         unit.Celsius,
		TempStep:                0.1,
	}
}

type Mode string

const (
	ModeAuto    Mode = "auto"
	ModeOff     Mode = "off"
	ModeCool    Mode = "cool"
	ModeHeat    Mode = "heat"
	ModeDry     Mode = "dry"
	ModeFanOnly Mode = "fan_only"
)

type FanMode string

const (
	FanModeAuto   FanMode = "auto"
	FanModeLow    FanMode = "low"
	FanModeMedium FanMode = "medium"
	FanModeHigh   FanMode = "high"
)

type PresetMode string

const (
	PresetEco      PresetMode = "eco"
	PresetAway     PresetMode = "away"
	PresetBoost    PresetMode = "boost"
	PresetComfort  PresetMode = "comfort"
	PresetHome     PresetMode = "home"
	PresetSleep    PresetMode = "sleep"
	PresetActivity PresetMode = "activity"
)

type SwingMode string

const (
	SwingModeOn  SwingMode = "on"
	SwingModeOff SwingMode = "off"
)

type SwingModeHorizontal string

const (
	SwingModeHorizontalOn  SwingModeHorizontal = "on"
	SwingModeHorizontalOff SwingModeHorizontal = "off"
)
