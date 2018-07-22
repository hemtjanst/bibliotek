package feature

type Info struct {
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	Step     int    `json:"step,omitempty"`
	GetTopic string `json:"getTopic,omitempty"`
	SetTopic string `json:"setTopic,omitempty"`
}

type Feature interface {
	Name() string
	Min() int
	Max() int
	Step() int
}

type FeatureTransporter interface {
	SubscribeFeature(*Info) chan []byte
	UpdateFeature(*Info, []byte)
	SetFeature(*Info, []byte)
}

type feature struct {
	info        *Info
	name        string
	transporter FeatureTransporter
}

func New(name string, info *Info, transporter FeatureTransporter) Feature {
	return &feature{
		name:        name,
		info:        info,
		transporter: transporter,
	}
}

func (f *feature) Name() string { return f.name }
func (f *feature) Min() int     { return f.info.Min }
func (f *feature) Max() int     { return f.info.Max }
func (f *feature) Step() int    { return f.info.Step }
