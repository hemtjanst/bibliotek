package device

// Origin is the origin information.
//
// It represents where a device is coming from.
type Origin struct {
	Name            string `json:"name,omitempty"`
	SoftwareVersion string `json:"sw,omitempty"`
	URL             string `json:"url,omitempty"`
}
