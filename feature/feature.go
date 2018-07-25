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
	Set(string) error
	Update(string) error
	OnSet() (chan string, error)
	OnUpdate() (chan string, error)
	Exists() bool
}

type Transport interface {
	SubscribeFeature(string) chan []byte
	UpdateFeature(*Info, []byte)
	SetFeature(*Info, []byte)
}

type feature struct {
	info      *Info
	name      string
	transport Transport
}

func New(name string, info *Info, transport Transport) Feature {
	return &feature{
		name:      name,
		info:      info,
		transport: transport,
	}
}

func (f *feature) Name() string { return f.name }
func (f *feature) Min() int     { return f.info.Min }
func (f *feature) Max() int     { return f.info.Max }
func (f *feature) Step() int    { return f.info.Step }
func (f *feature) Exists() bool { return true }

func (f *feature) Update(s string) error {
	f.transport.UpdateFeature(f.info, []byte(s))
	return nil
}

func (f *feature) Set(s string) error {
	f.transport.SetFeature(f.info, []byte(s))
	return nil
}

func (f *feature) OnUpdate() (chan string, error) {
	res := f.transport.SubscribeFeature(f.info.GetTopic)
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
	return ch, nil
}

func (f *feature) OnSet() (chan string, error) {
	res := f.transport.SubscribeFeature(f.info.SetTopic)
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
	return ch, nil
}
