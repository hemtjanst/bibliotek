package hass

import (
	"encoding/json"
	"fmt"
)

type Availability struct {
	PayloadAvailable    string `json:"payload_available,omitempty"`
	PayloadNotAvailable string `json:"payload_not_available,omitempty"`
	Topic               string `json:"topic,omitempty"`
	ValueTemplate       string `json:"value_template,omitempty"`
}

type ConnInfo struct {
	Type string
	ID   string
}

func (c ConnInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{c.Type, c.ID})
}

func (c *ConnInfo) UnmarshalJSON(b []byte) error {
	var d []string
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	if len(d) != 2 {
		return fmt.Errorf("connection info must be a string array of size 2")
	}
	c.Type, c.ID = d[0], d[1]
	return nil
}

type DeviceInfo struct {
	ConfigurationUrl string     `json:"configuration_url,omitempty"`
	Connections      []ConnInfo `json:"connections,omitempty"`
	Identifiers      []string   `json:"identifiers,omitempty"`
	Manufacturer     string     `json:"manufacturer,omitempty"`
	Model            string     `json:"model,omitempty"`
	Name             string     `json:"name,omitempty"`
	SwVersion        string     `json:"sw_version,omitempty"`
	SerialNumber     string     `json:"serial_number,omitempty"`
}

type Origin struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	SwVersion  string `json:"sw_version,omitempty" yaml:"sw_version,omitempty"`
	SupportUrl string `json:"support_url,omitempty" yaml:"support_url,omitempty"`
}

type Device struct {
	Device     *DeviceInfo           `json:"device,omitempty" yaml:"device,omitempty"`
	Origin     *Origin               `json:"origin,omitempty" yaml:"origin,omitempty"`
	Components map[string]*Component `json:"components,omitempty" yaml:"components,omitempty"`
	StateTopic string                `json:"state_topic,omitempty" yaml:"state_topic,omitempty"`
}

type Component struct {
	Platform                       string         `json:"p,omitempty" yaml:"p,omitempty"`
	BaseTopic                      string         `json:"~,omitempty"`
	ActionTopic                    string         `json:"act_t,omitempty"`
	ActionTemplate                 string         `json:"act_tpl,omitempty"`
	AutomationType                 string         `json:"atype,omitempty"`
	AuxCommandTopic                string         `json:"aux_cmd_t,omitempty"`
	AuxStateTemplate               string         `json:"aux_stat_tpl,omitempty"`
	AuxStateTopic                  string         `json:"aux_stat_t,omitempty"`
	AvailableTones                 string         `json:"av_tones,omitempty"`
	Availability                   []Availability `json:"avty,omitempty"`
	AvailabilityMode               string         `json:"avty_mode,omitempty"`
	AvailabilityTopic              string         `json:"avty_t,omitempty"`
	AvailabilityTemplate           string         `json:"avty_tpl,omitempty"`
	AwayModeCommandTopic           string         `json:"away_mode_cmd_t,omitempty"`
	AwayModeStateTemplate          string         `json:"away_mode_stat_tpl,omitempty"`
	AwayModeStateTopic             string         `json:"away_mode_stat_t,omitempty"`
	BlueTemplate                   string         `json:"b_tpl,omitempty"`
	BrightnessCommandTopic         string         `json:"bri_cmd_t,omitempty"`
	BrightnessCommandTemplate      string         `json:"bri_cmd_tpl,omitempty"`
	BrightnessScale                int            `json:"bri_scl,omitempty"`
	BrightnessStateTopic           string         `json:"bri_stat_t,omitempty"`
	BrightnessTemplate             string         `json:"bri_tpl,omitempty"`
	BrightnessValueTemplate        string         `json:"bri_val_tpl,omitempty"`
	ColorTempCommandTemplate       string         `json:"clr_temp_cmd_tpl,omitempty"`
	BatteryLevelTopic              string         `json:"bat_lev_t,omitempty"`
	BatteryLevelTemplate           string         `json:"bat_lev_tpl,omitempty"`
	ChargingTopic                  string         `json:"chrg_t,omitempty"`
	ChargingTemplate               string         `json:"chrg_tpl,omitempty"`
	ColorTempCommandTopic          string         `json:"clr_temp_cmd_t,omitempty"`
	ColorTempStateTopic            string         `json:"clr_temp_stat_t,omitempty"`
	ColorTempTemplate              string         `json:"clr_temp_tpl,omitempty"`
	ColorTempValueTemplate         string         `json:"clr_temp_val_tpl,omitempty"`
	CleaningTopic                  string         `json:"cln_t,omitempty"`
	CleaningTemplate               string         `json:"cln_tpl,omitempty"`
	CommandOffTemplate             string         `json:"cmd_off_tpl,omitempty"`
	CommandOnTemplate              string         `json:"cmd_on_tpl,omitempty"`
	CommandTopic                   string         `json:"cmd_t,omitempty"`
	CommandTemplate                string         `json:"cmd_tpl,omitempty"`
	CodeArmRequired                string         `json:"cod_arm_req,omitempty"`
	CodeDisarmRequired             string         `json:"cod_dis_req,omitempty"`
	CodeTriggerRequired            string         `json:"cod_trig_req,omitempty"`
	CurrentTemperatureTopic        string         `json:"curr_temp_t,omitempty"`
	CurrentTemperatureTemplate     string         `json:"curr_temp_tpl,omitempty"`
	Device                         *DeviceInfo    `json:"dev,omitempty"`
	DeviceClass                    string         `json:"dev_cla,omitempty"`
	DockedTopic                    string         `json:"dock_t,omitempty"`
	DockedTemplate                 string         `json:"dock_tpl,omitempty"`
	Encoding                       string         `json:"e,omitempty"`
	EntityCategory                 string         `json:"ent_cat,omitempty"`
	ErrorTopic                     string         `json:"err_t,omitempty"`
	ErrorTemplate                  string         `json:"err_tpl,omitempty"`
	FanSpeedTopic                  string         `json:"fanspd_t,omitempty"`
	FanSpeedTemplate               string         `json:"fanspd_tpl,omitempty"`
	FanSpeedList                   string         `json:"fanspd_lst,omitempty"`
	FlashTimeLong                  string         `json:"flsh_tlng,omitempty"`
	FlashTimeShort                 string         `json:"flsh_tsht,omitempty"`
	EffectCommandTopic             string         `json:"fx_cmd_t,omitempty"`
	EffectCommandTemplate          string         `json:"fx_cmd_tpl,omitempty"`
	EffectList                     string         `json:"fx_list,omitempty"`
	EffectStateTopic               string         `json:"fx_stat_t,omitempty"`
	EffectTemplate                 string         `json:"fx_tpl,omitempty"`
	EffectValueTemplate            string         `json:"fx_val_tpl,omitempty"`
	ExpireAfter                    string         `json:"exp_aft,omitempty"`
	FanModeCommandTemplate         string         `json:"fan_mode_cmd_tpl,omitempty"`
	FanModeCommandTopic            string         `json:"fan_mode_cmd_t,omitempty"`
	FanModeStateTemplate           string         `json:"fan_mode_stat_tpl,omitempty"`
	FanModeStateTopic              string         `json:"fan_mode_stat_t,omitempty"`
	ForceUpdate                    string         `json:"frc_upd,omitempty"`
	GreenTemplate                  string         `json:"g_tpl,omitempty"`
	HoldCommandTemplate            string         `json:"hold_cmd_tpl,omitempty"`
	HoldCommandTopic               string         `json:"hold_cmd_t,omitempty"`
	HoldStateTemplate              string         `json:"hold_stat_tpl,omitempty"`
	HoldStateTopic                 string         `json:"hold_stat_t,omitempty"`
	HsCommandTopic                 string         `json:"hs_cmd_t,omitempty"`
	HsStateTopic                   string         `json:"hs_stat_t,omitempty"`
	HsValueTemplate                string         `json:"hs_val_tpl,omitempty"`
	Icon                           string         `json:"ic,omitempty"`
	Initial                        string         `json:"init,omitempty"`
	TargetHumidityCommandTopic     string         `json:"hum_cmd_t,omitempty"`
	TargetHumidityCommandTemplate  string         `json:"hum_cmd_tpl,omitempty"`
	TargetHumidityStateTopic       string         `json:"hum_stat_t,omitempty"`
	TargetHumidityStateTemplate    string         `json:"hum_stat_tpl,omitempty"`
	JsonAttributes                 string         `json:"json_attr,omitempty"`
	JsonAttributesTopic            string         `json:"json_attr_t,omitempty"`
	JsonAttributesTemplate         string         `json:"json_attr_tpl,omitempty"`
	MaxMireds                      int            `json:"max_mirs,omitempty"`
	MinMireds                      int            `json:"min_mirs,omitempty"`
	MaxTemp                        string         `json:"max_temp,omitempty"`
	MinTemp                        string         `json:"min_temp,omitempty"`
	MaxHumidity                    string         `json:"max_hum,omitempty"`
	MinHumidity                    string         `json:"min_hum,omitempty"`
	ModeCommandTemplate            string         `json:"mode_cmd_tpl,omitempty"`
	ModeCommandTopic               string         `json:"mode_cmd_t,omitempty"`
	ModeStateTemplate              string         `json:"mode_stat_tpl,omitempty"`
	ModeStateTopic                 string         `json:"mode_stat_t,omitempty"`
	Modes                          string         `json:"modes,omitempty"`
	Name                           string         `json:"name,omitempty"`
	ObjectId                       string         `json:"obj_id,omitempty"`
	OffDelay                       string         `json:"off_dly,omitempty"`
	OnCommandType                  string         `json:"on_cmd_type,omitempty"`
	Optimistic                     string         `json:"opt,omitempty"`
	OscillationCommandTopic        string         `json:"osc_cmd_t,omitempty"`
	OscillationCommandTemplate     string         `json:"osc_cmd_tpl,omitempty"`
	OscillationStateTopic          string         `json:"osc_stat_t,omitempty"`
	OscillationValueTemplate       string         `json:"osc_val_tpl,omitempty"`
	PercentageCommandTopic         string         `json:"pct_cmd_t,omitempty"`
	PercentageCommandTemplate      string         `json:"pct_cmd_tpl,omitempty"`
	PercentageStateTopic           string         `json:"pct_stat_t,omitempty"`
	PercentageValueTemplate        string         `json:"pct_val_tpl,omitempty"`
	Payload                        string         `json:"pl,omitempty"`
	PayloadArmAway                 string         `json:"pl_arm_away,omitempty"`
	PayloadArmHome                 string         `json:"pl_arm_home,omitempty"`
	PayloadArmCustomBypass         string         `json:"pl_arm_custom_b,omitempty"`
	PayloadArmNight                string         `json:"pl_arm_nite,omitempty"`
	PayloadAvailable               string         `json:"pl_avail,omitempty"`
	PayloadCleanSpot               string         `json:"pl_cln_sp,omitempty"`
	PayloadClose                   string         `json:"pl_cls,omitempty"`
	PayloadDisarm                  string         `json:"pl_disarm,omitempty"`
	PayloadHome                    string         `json:"pl_home,omitempty"`
	PayloadLock                    string         `json:"pl_lock,omitempty"`
	PayloadLocate                  string         `json:"pl_loc,omitempty"`
	PayloadNotAvailable            string         `json:"pl_not_avail,omitempty"`
	PayloadNotHome                 string         `json:"pl_not_home,omitempty"`
	PayloadOff                     string         `json:"pl_off,omitempty"`
	PayloadOn                      string         `json:"pl_on,omitempty"`
	PayloadOpen                    string         `json:"pl_open,omitempty"`
	PayloadOscillationOff          string         `json:"pl_osc_off,omitempty"`
	PayloadOscillationOn           string         `json:"pl_osc_on,omitempty"`
	PayloadPause                   string         `json:"pl_paus,omitempty"`
	PayloadStop                    string         `json:"pl_stop,omitempty"`
	PayloadStart                   string         `json:"pl_strt,omitempty"`
	PayloadStartPause              string         `json:"pl_stpa,omitempty"`
	PayloadReturnToBase            string         `json:"pl_ret,omitempty"`
	PayloadResetHumidity           string         `json:"pl_rst_hum,omitempty"`
	PayloadResetMode               string         `json:"pl_rst_mode,omitempty"`
	PayloadResetPercentage         string         `json:"pl_rst_pct,omitempty"`
	PayloadResetPresetMode         string         `json:"pl_rst_pr_mode,omitempty"`
	PayloadTurnOff                 string         `json:"pl_toff,omitempty"`
	PayloadTurnOn                  string         `json:"pl_ton,omitempty"`
	PayloadTrigger                 string         `json:"pl_trig,omitempty"`
	PayloadUnlock                  string         `json:"pl_unlk,omitempty"`
	PositionClosed                 int            `json:"pos_clsd,omitempty"`
	PositionOpen                   int            `json:"pos_open,omitempty"`
	PowerCommandTopic              string         `json:"pow_cmd_t,omitempty"`
	PowerStateTopic                string         `json:"pow_stat_t,omitempty"`
	PowerStateTemplate             string         `json:"pow_stat_tpl,omitempty"`
	PresetModeCommandTopic         string         `json:"pr_mode_cmd_t,omitempty"`
	PresetModeCommandTemplate      string         `json:"pr_mode_cmd_tpl,omitempty"`
	PresetModeStateTopic           string         `json:"pr_mode_stat_t,omitempty"`
	PresetModeValueTemplate        string         `json:"pr_mode_val_tpl,omitempty"`
	PresetModes                    string         `json:"pr_modes,omitempty"`
	RedTemplate                    string         `json:"r_tpl,omitempty"`
	Retain                         string         `json:"ret,omitempty"`
	RgbCommandTemplate             string         `json:"rgb_cmd_tpl,omitempty"`
	RgbCommandTopic                string         `json:"rgb_cmd_t,omitempty"`
	RgbStateTopic                  string         `json:"rgb_stat_t,omitempty"`
	RgbValueTemplate               string         `json:"rgb_val_tpl,omitempty"`
	SendCommandTopic               string         `json:"send_cmd_t,omitempty"`
	SendIfOff                      string         `json:"send_if_off,omitempty"`
	SetFanSpeedTopic               string         `json:"set_fan_spd_t,omitempty"`
	SetPositionTemplate            string         `json:"set_pos_tpl,omitempty"`
	SetPositionTopic               string         `json:"set_pos_t,omitempty"`
	PositionTopic                  string         `json:"pos_t,omitempty"`
	PositionTemplate               string         `json:"pos_tpl,omitempty"`
	SpeedRangeMin                  string         `json:"spd_rng_min,omitempty"`
	SpeedRangeMax                  string         `json:"spd_rng_max,omitempty"`
	SourceType                     string         `json:"src_type,omitempty"`
	StateClass                     string         `json:"stat_cla,omitempty"`
	StateClosed                    string         `json:"stat_clsd,omitempty"`
	StateClosing                   string         `json:"stat_closing,omitempty"`
	StateOff                       string         `json:"stat_off,omitempty"`
	StateOn                        string         `json:"stat_on,omitempty"`
	StateOpen                      string         `json:"stat_open,omitempty"`
	StateOpening                   string         `json:"stat_opening,omitempty"`
	StateStopped                   string         `json:"stat_stopped,omitempty"`
	StateLocked                    string         `json:"stat_locked,omitempty"`
	StateUnlocked                  string         `json:"stat_unlocked,omitempty"`
	StateTopic                     string         `json:"stat_t,omitempty"`
	StateTemplate                  string         `json:"stat_tpl,omitempty"`
	StateValueTemplate             string         `json:"stat_val_tpl,omitempty"`
	Subtype                        string         `json:"stype,omitempty"`
	SupportDuration                string         `json:"sup_duration,omitempty"`
	SupportVolumeSet               string         `json:"sup_vol,omitempty"`
	SupportedFeatures              string         `json:"sup_feat,omitempty"`
	SupportedTurnOff               string         `json:"sup_off,omitempty"`
	SwingModeCommandTemplate       string         `json:"swing_mode_cmd_tpl,omitempty"`
	SwingModeCommandTopic          string         `json:"swing_mode_cmd_t,omitempty"`
	SwingModeStateTemplate         string         `json:"swing_mode_stat_tpl,omitempty"`
	SwingModeStateTopic            string         `json:"swing_mode_stat_t,omitempty"`
	TemperatureCommandTemplate     string         `json:"temp_cmd_tpl,omitempty"`
	TemperatureCommandTopic        string         `json:"temp_cmd_t,omitempty"`
	TemperatureHighCommandTemplate string         `json:"temp_hi_cmd_tpl,omitempty"`
	TemperatureHighCommandTopic    string         `json:"temp_hi_cmd_t,omitempty"`
	TemperatureHighStateTemplate   string         `json:"temp_hi_stat_tpl,omitempty"`
	TemperatureHighStateTopic      string         `json:"temp_hi_stat_t,omitempty"`
	TemperatureLowCommandTemplate  string         `json:"temp_lo_cmd_tpl,omitempty"`
	TemperatureLowCommandTopic     string         `json:"temp_lo_cmd_t,omitempty"`
	TemperatureLowStateTemplate    string         `json:"temp_lo_stat_tpl,omitempty"`
	TemperatureLowStateTopic       string         `json:"temp_lo_stat_t,omitempty"`
	TemperatureStateTemplate       string         `json:"temp_stat_tpl,omitempty"`
	TemperatureStateTopic          string         `json:"temp_stat_t,omitempty"`
	TemperatureUnit                string         `json:"temp_unit,omitempty"`
	TiltClosedValue                string         `json:"tilt_clsd_val,omitempty"`
	TiltCommandTopic               string         `json:"tilt_cmd_t,omitempty"`
	TiltCommandTemplate            string         `json:"tilt_cmd_tpl,omitempty"`
	TiltInvertState                string         `json:"tilt_inv_stat,omitempty"`
	TiltMax                        string         `json:"tilt_max,omitempty"`
	TiltMin                        string         `json:"tilt_min,omitempty"`
	TiltOpenedValue                string         `json:"tilt_opnd_val,omitempty"`
	TiltOptimistic                 string         `json:"tilt_opt,omitempty"`
	TiltStatusTopic                string         `json:"tilt_status_t,omitempty"`
	TiltStatusTemplate             string         `json:"tilt_status_tpl,omitempty"`
	Topic                          string         `json:"t,omitempty"`
	UniqueId                       string         `json:"uniq_id,omitempty"`
	UnitOfMeasurement              string         `json:"unit_of_meas,omitempty"`
	ValueTemplate                  string         `json:"val_tpl,omitempty"`
	WhiteValueCommandTopic         string         `json:"whit_val_cmd_t,omitempty"`
	WhiteValueScale                string         `json:"whit_val_scl,omitempty"`
	WhiteValueStateTopic           string         `json:"whit_val_stat_t,omitempty"`
	WhiteValueTemplate             string         `json:"whit_val_tpl,omitempty"`
	XyCommandTopic                 string         `json:"xy_cmd_t,omitempty"`
	XyStateTopic                   string         `json:"xy_stat_t,omitempty"`
	XyValueTemplate                string         `json:"xy_val_tpl,omitempty"`
}
