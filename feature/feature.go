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
	Set(string)
	Update(string)
	OnSet() chan string
	OnUpdate() chan string
}

type FeatureTransporter interface {
	SubscribeFeature(string) chan []byte
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

func (f *feature) Update(s string) {
	f.transporter.UpdateFeature(f.info, []byte(s))
}

func (f *feature) Set(s string) {
	f.transporter.SetFeature(f.info, []byte(s))
}

func (f *feature) OnUpdate() chan string {
	res := f.transporter.SubscribeFeature(f.info.GetTopic)
	ch := make(chan string, 5)
	go func() {
		for {
			msg, open := <-res
			if !open {
				close(ch)
				return
			}
			ch <- string(msg)
		}
	}()
	return ch
}

func (f *feature) OnSet() chan string {
	res := f.transporter.SubscribeFeature(f.info.SetTopic)
	ch := make(chan string, 5)
	go func() {
		for {
			msg, open := <-res
			if !open {
				close(ch)
				return
			}
			ch <- string(msg)
		}
	}()
	return ch
}
