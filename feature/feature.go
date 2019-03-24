package feature

import "strconv"

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
	Exists() bool
	Set(string) error
	OnUpdate() (chan string, error)
	Update(string) error
	OnSet() (chan string, error)
	UpdateInfo(*Info) []*InfoUpdate
	GetTopic() string
	SetTopic() string
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

type InfoUpdate struct {
	Name string
	Old  string
	New  string
}

func New(name string, info *Info, transport Transport) Feature {
	return &feature{
		name:      name,
		info:      info,
		transport: transport,
	}
}

func (f *feature) Name() string     { return f.name }
func (f *feature) Min() int         { return f.info.Min }
func (f *feature) Max() int         { return f.info.Max }
func (f *feature) Step() int        { return f.info.Step }
func (f *feature) Exists() bool     { return true }
func (f *feature) GetTopic() string { return f.info.GetTopic }
func (f *feature) SetTopic() string { return f.info.SetTopic }

func (f *feature) UpdateInfo(i *Info) (u []*InfoUpdate) {
	if i.GetTopic != "" && i.GetTopic != f.info.GetTopic {
		u = append(u, &InfoUpdate{"getTopic", f.info.GetTopic, i.GetTopic})
		f.info.GetTopic = i.GetTopic
	}
	if i.SetTopic != "" && i.SetTopic != f.info.SetTopic {
		u = append(u, &InfoUpdate{"setTopic", f.info.SetTopic, i.SetTopic})
		f.info.SetTopic = i.SetTopic
	}
	if i.Min != f.info.Min {
		u = append(u, &InfoUpdate{"min", strconv.Itoa(f.info.Min), strconv.Itoa(i.Min)})
		f.info.Min = i.Min
	}
	if i.Max != f.info.Max {
		u = append(u, &InfoUpdate{"max", strconv.Itoa(f.info.Max), strconv.Itoa(i.Max)})
		f.info.Max = i.Max
	}
	if i.Step != f.info.Step {
		u = append(u, &InfoUpdate{"step", strconv.Itoa(f.info.Step), strconv.Itoa(i.Step)})
		f.info.Step = i.Step
	}
	return
}

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
