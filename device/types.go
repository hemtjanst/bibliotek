package device

import "strings"

type Type string

func (t Type) Equal(y Type) bool {
	return strings.ToLower(string(t)) == strings.ToLower(string(y))
}

const (
	AccessoryInformation         Type = "accessoryInformation"
	AirPurifier                  Type = "airPurifier"
	AirQualitySensor             Type = "airQualitySensor"
	BatteryService               Type = "batteryService"
	BridgeConfiguration          Type = "bridgeConfiguration"
	BridgingState                Type = "bridgingState"
	CameraControl                Type = "cameraControl"
	CameraRTPStreamManagement    Type = "cameraRTPStreamManagement"
	CarbonDioxideSensor          Type = "carbonDioxideSensor"
	CarbonMonoxideSensor         Type = "carbonMonoxideSensor"
	ContactSensor                Type = "contactSensor"
	Door                         Type = "door"
	Doorbell                     Type = "doorbell"
	Fan                          Type = "fan"
	FanV2                        Type = "fanV2"
	Faucet                       Type = "faucet"
	FilterMaintenance            Type = "filterMaintenance"
	GarageDoorOpener             Type = "garageDoorOpener"
	HeaterCooler                 Type = "heaterCooler"
	HumidifierDehumidifier       Type = "humidifierDehumidifier"
	HumiditySensor               Type = "humiditySensor"
	IrrigationSystem             Type = "irrigationSystem"
	LeakSensor                   Type = "leakSensor"
	LightSensor                  Type = "lightSensor"
	Lightbulb                    Type = "lightbulb"
	LockManagement               Type = "lockManagement"
	LockMechanism                Type = "lockMechanism"
	Microphone                   Type = "microphone"
	MotionSensor                 Type = "motionSensor"
	OccupancySensor              Type = "occupancySensor"
	Outlet                       Type = "outlet"
	SecuritySystem               Type = "securitySystem"
	ServiceLabel                 Type = "serviceLabel"
	Slat                         Type = "slat"
	SmokeSensor                  Type = "smokeSensor"
	Speaker                      Type = "speaker"
	StatefulProgrammableSwitch   Type = "statefulProgrammableSwitch"
	StatelessProgrammableSwitch  Type = "statelessProgrammableSwitch"
	Switch                       Type = "switch"
	TemperatureSensor            Type = "temperatureSensor"
	Thermostat                   Type = "thermostat"
	TimeInformation              Type = "timeInformation"
	TunneledBTLEAccessoryService Type = "tunneledBTLEAccessoryService"
	Valve                        Type = "valve"
	Window                       Type = "window"
	WindowCovering               Type = "windowCovering"
)
