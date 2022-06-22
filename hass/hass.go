package hass

import (
	"fmt"
	"lib.hemtjan.st/device"
	"lib.hemtjan.st/feature"
	"strings"
)

func HTtoHA(dev *device.Info) (out []*HassDev, err error) {

	withFt := func(name feature.Type, cb func(ft *feature.Info, get, set string)) bool {
		if ft, ok := dev.Features[string(name)]; ok {
			get := fmt.Sprintf("~/%s/get", name)
			set := fmt.Sprintf("~/%s/set", name)
			if ft.GetTopic != "" {
				get = ft.GetTopic
			}
			if ft.SetTopic != "" {
				set = ft.SetTopic
			}
			cb(ft, get, set)
			return true
		}
		return false
	}

	o := &HassDev{}
	out = append(out, o)

	nameParts := strings.Split(dev.Topic, "")
	for i, p := range nameParts {
		p = strings.ToLower(p)
		if !(p >= "a" && p <= "z" || p >= "0" && p <= "9" || p == "_") {
			p = "_"
		}
		nameParts[i] = p
	}

	onFt := func() bool {
		return withFt(feature.On, func(ft *feature.Info, get, set string) {
			o.StateOn = "1"
			o.StateOff = "0"
			o.PayloadOn = "1"
			o.PayloadOff = "0"
			o.StateTopic = get
			o.CommandTopic = set
		})
	}

	uniq := strings.Join(nameParts, "")
	o.BaseTopic = dev.Topic
	o.UniqueId = "mqtt_" + uniq
	o.Name = dev.Name
	o.Device = &HassDevInfo{
		Identifiers:  []string{o.UniqueId},
		Manufacturer: dev.Manufacturer,
		Model:        dev.Model,
		Name:         dev.Name,
	}

	withFt(feature.BatteryLevel, func(ft *feature.Info, get, set string) { o.BatteryLevelTopic = get })
	withFt(feature.ChargingState, func(ft *feature.Info, get, set string) { o.ChargingTopic = get })

	switch device.Type(dev.Type) {
	case device.Lightbulb:
		o.Type = "light"
		o.DeviceClass = "light"

		if !onFt() {
			return nil, fmt.Errorf("light requires feature on")
		}

		withFt(feature.Brightness, func(ft *feature.Info, get, set string) {
			o.BrightnessStateTopic = get
			o.BrightnessCommandTopic = set
			o.BrightnessScale = 100
		})

		withFt(feature.ColorTemperature, func(ft *feature.Info, get, set string) {
			o.ColorTempStateTopic = get
			o.ColorTempCommandTopic = set
			if ft.Min > 0 {
				o.MinMireds = ft.Min
			}
			if ft.Max > 0 {
				o.MaxMireds = ft.Max
			}
		})

		withFt(feature.Color, func(ft *feature.Info, get, set string) {
			o.RgbStateTopic = get
			o.RgbCommandTopic = set
		})

	case device.Outlet, device.Switch:
		o.Type = "switch"
		o.DeviceClass = dev.Type
		if !onFt() {
			return nil, fmt.Errorf("outlet/switch requires feature on")
		}

		withFt(feature.CurrentPower, func(ft *feature.Info, get, set string) {

		})

		withFt(feature.EnergyUsed, func(ft *feature.Info, get, set string) {

		})

	case device.WindowCovering:
		o.Type = "cover"
		o.DeviceClass = "shade"
		if !withFt(feature.CurrentPosition, func(ft *feature.Info, get, set string) {
			o.PositionOpen = 100
			o.PositionClosed = 0
			o.PositionTopic = get
		}) || !withFt(feature.TargetPosition, func(ft *feature.Info, get, set string) {
			o.PayloadClose = "0"
			o.PayloadOpen = "100"
			o.PayloadStop = "-1"
			o.SetPositionTopic = set
			o.CommandTopic = set
		}) || !withFt(feature.PositionState, func(ft *feature.Info, get, set string) {
			o.StateTopic = get
			o.StateOpening = "1"
			o.StateClosing = "0"
			o.StateStopped = "2"
		}) {
			return nil, fmt.Errorf("window covering requires features currentPosition, targetPosition and positionState")
		}
	case device.TemperatureSensor:
		o.Type = "sensor"
		o.DeviceClass = "temperature"
		o.UnitOfMeasurement = "°C"
		if !withFt(feature.CurrentTemperature, func(ft *feature.Info, get, set string) {
			o.StateTopic = get
		}) {
			return nil, fmt.Errorf("temperature sensor requires feature currentTemperature")
		}
	case device.HumiditySensor:
		o.Type = "sensor"
		o.DeviceClass = "humidity"
		o.UnitOfMeasurement = "%"
		if !withFt(feature.CurrentRelativeHumidity, func(ft *feature.Info, get, set string) {
			o.StateTopic = get
		}) {
			return nil, fmt.Errorf("humidity sensor requires feature currentRelativeHumidity")
		}
	case device.WeatherStation:
		out = []*HassDev{}
		for n, ft := range dev.Features {
			d := &HassDev{
				Device: &HassDevInfo{
					Identifiers:  []string{o.UniqueId},
					Manufacturer: dev.Manufacturer,
					Model:        dev.Model,
					Name:         dev.Name,
				},
				Type:       "sensor",
				BaseTopic:  o.BaseTopic + "/" + n,
				StateTopic: "~/get",
			}
			if ft.GetTopic != "" {
				d.StateTopic = ft.GetTopic
			}
			switch feature.Type(n) {
			case feature.CurrentRelativeHumidity:
				d.UniqueId = o.UniqueId + "_humidity"
				d.Name = o.Name + " Humidity"
				d.UnitOfMeasurement = "%"
				d.DeviceClass = "humidity"
			case feature.CurrentTemperature:
				d.UniqueId = o.UniqueId + "_temp"
				d.Name = o.Name + " Temperature"
				d.UnitOfMeasurement = "°C"
				d.DeviceClass = "temperature"
			case feature.Type("airPressure"):
				d.UniqueId = o.UniqueId + "_pressure"
				d.Name = o.Name + " Air Pressure"
				d.UnitOfMeasurement = "hPa"
				d.DeviceClass = "pressure"
			default:
				continue
			}
			out = append(out, d)
		}
	default:
		return nil, fmt.Errorf("unsupported type: %s", dev.Type)
	}

	return
}
