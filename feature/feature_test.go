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
	assert.True(t, f.Exists())
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
	err := f.Update("test1")
	assert.Nil(t, err)
	d.On("SetFeature", info, []byte("test2")).Return()
	err = f.Set("test2")
	assert.Nil(t, err)
	d.AssertExpectations(t)

	d.On("SubscribeFeature", "test/get").Return()
	res, err := f.OnUpdate()
	assert.Nil(t, err)
	d.subChan <- []byte("test3")
	msg := <-res
	assert.Equal(t, "test3", msg)
	close(d.subChan)
	_, open := <-res
	assert.False(t, open)
	d.AssertExpectations(t)

	d.On("SubscribeFeature", "test/set").Return()
	res, err = f.OnSet()
	assert.Nil(t, err)
	d.subChan <- []byte("test4")
	msg = <-res
	assert.Equal(t, "test4", msg)
	close(d.subChan)
	_, open = <-res
	assert.False(t, open)
	d.AssertExpectations(t)
}
