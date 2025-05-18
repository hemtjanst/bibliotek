package component

type StateClass string

const (
	Measurement      StateClass = "measurement"
	MeasurementAngle StateClass = "measurement_angle"
	Total            StateClass = "total"
	TotalIncreasing  StateClass = "total_increasing"
)
