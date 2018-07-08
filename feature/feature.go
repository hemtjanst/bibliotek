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

type feature struct {
	info *Info
	name string
}

func New(name string, info *Info) Feature {
	return &feature{
		name: name,
		info: info,
	}
}

func (f *feature) Name() string { return f.name }
func (f *feature) Min() int     { return f.info.Min }
func (f *feature) Max() int     { return f.info.Max }
func (f *feature) Step() int    { return f.info.Step }
