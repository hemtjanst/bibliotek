package feature // import "lib.hemtjan.st/feature"

import (
	"strconv"
	"time"
)

// Info holds information about a feature
type Info struct {
	Min      int    `json:"min,omitempty"`
	Max      int    `json:"max,omitempty"`
	Step     int    `json:"step,omitempty"`
	GetTopic string `json:"getTopic,omitempty"`
	SetTopic string `json:"setTopic,omitempty"`
}

// Feature represents a feature of a device
type Feature interface {
	Name() string
	Min() int
	Max() int
	Step() int
	Exists() bool
	Set(string) error
	OnUpdate() (chan string, error)
	OnUpdateFunc(func(string)) error
	Update(string) error
	OnSet() (chan string, error)
	OnSetFunc(func(string)) error
	Value() string
	UpdateInfo(*Info) []*InfoUpdate
	GetTopic() string
	SetTopic() string
}

// Transport is the feature's transport
type Transport interface {
	SubscribeFeature(string) chan []byte
	UpdateFeature(*Info, []byte)
	SetFeature(*Info, []byte)
}

type feature struct {
	info      *Info
	name      string
	transport Transport
	value     string
	updateSub bool
}

// InfoUpdate represents an update to the feature's current
// info
type InfoUpdate struct {
	Name string
	Old  string
	New  string
}

// New creates a new feature with the specified name and info
// and embeds the transport over which it can be interacted with
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

// UpdateInfo updates the Info of the feature
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

// Update publishes the new value of the feature
// This is a no-op if the new and current value are the same
func (f *feature) Update(s string) error {
	if f.value == s {
		return nil
	}
	f.transport.UpdateFeature(f.info, []byte(s))
	f.value = s
	return nil
}

// Set updates the value of the feature
func (f *feature) Set(s string) error {
	f.transport.SetFeature(f.info, []byte(s))
	return nil
}

// OnUpdate returns a channel on which updates
// of the feature's value are published
func (f *feature) OnUpdate() (chan string, error) {
	res := f.transport.SubscribeFeature(f.info.GetTopic)
	ch := make(chan string, 5)
	go func() {
		var value string
		for {
			msg, open := <-res
			if !open {
				close(ch)
				return
			}
			smsg := string(msg)
			if value == smsg {
				continue
			}
			value = smsg
			ch <- smsg
		}
	}()
	return ch, nil
}

// OnUpdateFunc calls the provided callback when updates
// of the feature's value are published
func (f *feature) OnUpdateFunc(fn func(val string)) error {
	ch, err := f.OnUpdate()
	if err != nil {
		return err
	}
	go func() {
		for {
			val, open := <-ch
			if !open {
				return
			}
			fn(val)
		}
	}()
	return nil
}

func (f *feature) Value() string {
	if !f.updateSub {
		f.updateSub = true
		ch := make(chan struct{})
		go func(ch chan struct{}) {
			defer func() {
				if ch != nil {
					close(ch)
				}
			}()
			up, err := f.OnUpdate()
			if err != nil {
				return
			}
			for {
				val, open := <-up
				if !open {
					f.updateSub = false
					return
				}
				f.value = string(val)
				if ch != nil {
					close(ch)
					ch = nil
				}
			}
		}(ch)
		select {
		case <-ch:
		case <-time.After(1 * time.Second):
		}
	}
	return f.value
}

// OnSet returns a channel on which notifications of a
// feature's intended new value are published
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

// OnSetFunc calls the provided callback when a
// feature's intended new value is published
func (f *feature) OnSetFunc(fn func(val string)) error {
	ch, err := f.OnSet()
	if err != nil {
		return err
	}
	go func() {
		for {
			val, open := <-ch
			if !open {
				return
			}
			fn(val)
		}
	}()
	return nil
}
