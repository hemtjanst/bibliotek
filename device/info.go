package device

// Info is the device information.
type Info struct {
	ID              string   `json:"ids"`
	Name            string   `json:"name"`
	Manufacturer    string   `json:"mf,omitempty"`
	Model           string   `json:"mdl,omitempty"`
	ModelID         string   `json:"mdl_id,omitempty"`
	SoftwareVersion string   `json:"sw,omitempty"`
	SerialNumber    string   `json:"sn,omitempty"`
	HardwareVersion string   `json:"hw,omitempty"`
	SuggestedArea   string   `json:"sa,omitempty"`
	Connections     []string `json:"cns,omitempty"`
	URL             string   `json:"cu,omitempty"`
}
