package state

type Class string

const (
	Measurement      Class = "measurement"
	MeasurementAngle Class = "measurement_angle"
	Total            Class = "total"
	TotalIncreasing  Class = "total_increasing"
)
