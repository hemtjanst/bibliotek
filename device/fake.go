package device

import "github.com/hemtjanst/bibliotek/feature"

type Fake struct {
	Err   error
	Topic string
}

func (f *Fake) Id() string {
	return f.Topic
}
func (f *Fake) Name() string {
	return "fake"
}
func (f *Fake) Manufacturer() string {
	return "fake"
}
func (f *Fake) Model() string {
	return "fake"
}
func (f *Fake) SerialNumber() string {
	return "fake"
}
func (f *Fake) Type() string {
	return "fake"
}
func (f *Fake) Feature(name string) feature.Server {
	return &feature.Fake{Err: f.Err, FeatureName: name}
}
func (f *Fake) Exists() bool {
	return false
}
func (f *Fake) Features() []feature.Server {
	return []feature.Server{}
}
