package device

type DeviceInfo struct {
	Topic        string                  `json:"topic"`
	Name         string                  `json:"name"`
	Manufacturer string                  `json:"manufacturer"`
	Model        string                  `json:"model"`
	SerialNumber string                  `json:"serialNumber"`
	Type         string                  `json:"type"`
	LastWillID   string                  `json:"lastWillID,omitempty"`
	Features     map[string]*FeatureInfo `json:"feature"`
	Reachable    bool                    `json:"-"`
}

type FeatureInfo struct {
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	Step     int    `json:"step,omitempty"`
	GetTopic string `json:"getTopic,omitempty"`
	SetTopic string `json:"setTopic,omitempty"`
}

type Device struct {
	DeviceInfo
}

func New(info DeviceInfo) *Device {
	return &Device{info}
}
