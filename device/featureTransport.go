package device

import "github.com/hemtjanst/bibliotek/feature"

func (d *Device) SubscribeFeature(topic string) chan []byte {
	return d.Transport.Subscribe(topic)
}
func (d *Device) UpdateFeature(f *feature.Info, b []byte) {
	d.Transport.Publish(f.GetTopic, b, true)
}
func (d *Device) SetFeature(f *feature.Info, b []byte) {
	d.Transport.Publish(f.SetTopic, b, false)
}
