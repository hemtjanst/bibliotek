package unit

type Measurement string

const (
	Ampere      Measurement = "A"
	MilliAmpere Measurement = "mA"

	Celsius Measurement = "°C"
	Kelvin  Measurement = "K"

	Calorie     Measurement = "cal"
	KiloCalorie Measurement = "kcal"
	MegaCalorie Measurement = "Mcal"
	GigaCalorie Measurement = "Gcal"

	Joule     Measurement = "J"
	KiloJoule Measurement = "kJ"
	MegaJoule Measurement = "MJ"

	Lux Measurement = "lx"

	MicrogramsPerCubicMeter = "µg/m³"

	Millimeter Measurement = "mm"
	Centimeter Measurement = "cm"

	MilliWattHour Measurement = "mWh"
	WattHour      Measurement = "Wh"
	KiloWattHour  Measurement = "kWh"
	MegaWattHour  Measurement = "MWh"
	GigaWattHour  Measurement = "GWh"

	MicroVolt Measurement = "µV"
	MilliVolt Measurement = "mV"
	Volt      Measurement = "V"
	KiloVolt  Measurement = "kV"
	MegaVolt  Measurement = "MV"

	Percent Measurement = "%"
)
