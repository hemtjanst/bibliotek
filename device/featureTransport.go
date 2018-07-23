package device

import "github.com/hemtjanst/bibliotek/feature"

func (d *device) SubscribeFeature(f *feature.Info) chan []byte {
	return d.transporter.Subscribe(f.GetTopic)
}
func (d *device) UpdateFeature(f *feature.Info, b []byte) {
	d.transporter.Publish(f.GetTopic, b, true)
}
func (d *device) SetFeature(f *feature.Info, b []byte) {
	d.transporter.Publish(f.SetTopic, b, false)
}
