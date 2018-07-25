package device

import "github.com/hemtjanst/bibliotek/feature"

func (d *device) SubscribeFeature(topic string) chan []byte {
	return d.transport.Subscribe(topic)
}
func (d *device) UpdateFeature(f *feature.Info, b []byte) {
	d.transport.Publish(f.GetTopic, b, true)
}
func (d *device) SetFeature(f *feature.Info, b []byte) {
	d.transport.Publish(f.SetTopic, b, false)
}
