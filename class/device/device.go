package device

type Class string

const (
	Battery          Class = "battery"
	Current          Class = "current"
	Energy           Class = "energy"
	Illuminance      Class = "illuminance"
	Heat             Class = "heat"
	RelativeHumidity Class = "humidity"
	PM1              Class = "pm1"
	PM25             Class = "pm25"
	PM10             Class = "pm10"
	Power            Class = "power"
	Precipitation    Class = "precipitation"
	Temperature      Class = "temperature"
	Voltage          Class = "voltage"
)
