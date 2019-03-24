package server

import (
	"github.com/hemtjanst/bibliotek/device"
	"github.com/hemtjanst/bibliotek/feature"
)

type FakeDevice struct {
	Err   error
	Topic string
}

func (f *FakeDevice) Id() string {
	return f.Topic
}
func (f *FakeDevice) Name() string {
	return "fake"
}
func (f *FakeDevice) Manufacturer() string {
	return "fake"
}
func (f *FakeDevice) Model() string {
	return "fake"
}
func (f *FakeDevice) SerialNumber() string {
	return "fake"
}
func (f *FakeDevice) Type() string {
	return "fake"
}
func (f *FakeDevice) Feature(name string) Feature {
	return &feature.Fake{Err: f.Err, FeatureName: name}
}
func (f *FakeDevice) Exists() bool {
	return false
}
func (f *FakeDevice) Features() []Feature {
	return []Feature{}
}
func (f *FakeDevice) setReachability(r bool) {
	return
}
func (f *FakeDevice) update(*device.Info) {
	return
}
