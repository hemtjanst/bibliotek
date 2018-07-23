package device

import "github.com/hemtjanst/bibliotek/feature"

func (d *device) SubscribeFeature(topic string) chan []byte {
	return d.transporter.Subscribe(topic)
}
func (d *device) UpdateFeature(f *feature.Info, b []byte) {
	d.transporter.Publish(f.GetTopic, b, true)
}
func (d *device) SetFeature(f *feature.Info, b []byte) {
	d.transporter.Publish(f.SetTopic, b, false)
}
