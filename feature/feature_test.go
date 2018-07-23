package feature

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDevice struct {
	subChan chan []byte
	mock.Mock
}

func (d *mockDevice) SubscribeFeature(topic string) chan []byte {
	d.Called(topic)
	ch := make(chan []byte, 5)
	d.subChan = ch
	return ch
}
func (d *mockDevice) UpdateFeature(f *Info, b []byte) {
	d.Called(f, b)
}
func (d *mockDevice) SetFeature(f *Info, b []byte) {
	d.Called(f, b)
}

func TestCreateFeature(t *testing.T) {
	info := &Info{
		Min:      0,
		Max:      10,
		Step:     1,
		GetTopic: "test/get",
		SetTopic: "test/set",
	}

	f := New("on", info, nil)
	assert.Equal(t, "on", f.Name())
	assert.Equal(t, 0, f.Min())
	assert.Equal(t, 10, f.Max())
	assert.Equal(t, 1, f.Step())
}

func TestFeatureTransporter(t *testing.T) {
	d := &mockDevice{}
	info := &Info{
		Min:      0,
		Max:      10,
		Step:     1,
		GetTopic: "test/get",
		SetTopic: "test/set",
	}

	f := New("on", info, d)
	d.On("UpdateFeature", info, []byte("test1")).Return()
	f.Update("test1")
	d.On("SetFeature", info, []byte("test2")).Return()
	f.Set("test2")
	d.AssertExpectations(t)

	d.On("SubscribeFeature", "test/get").Return()
	res := f.OnUpdate()
	d.subChan <- []byte("test3")
	msg := <-res
	assert.Equal(t, "test3", msg)
	close(d.subChan)
	_, open := <-res
	assert.False(t, open)
	d.AssertExpectations(t)

	d.On("SubscribeFeature", "test/set").Return()
	res = f.OnSet()
	d.subChan <- []byte("test4")
	msg = <-res
	assert.Equal(t, "test4", msg)
	close(d.subChan)
	_, open = <-res
	assert.False(t, open)
	d.AssertExpectations(t)
}
