package feature

import "testing"

func TestCreateFeature(t *testing.T) {
	info := &Info{
		Min:      0,
		Max:      10,
		Step:     1,
		GetTopic: "test/get",
		SetTopic: "test/set",
	}

	f := New("on", info)

	if f.Name() != "on" {
		t.Errorf("Expected: on, got: %s", f.Name())
	}
	if f.Min() != 0 {
		t.Errorf("Expected: 0, got: %d", f.Min())
	}
	if f.Max() != 10 {
		t.Errorf("Expected: 10, got: %d", f.Max())
	}
	if f.Step() != 1 {
		t.Errorf("Expected: 1, got: %d", f.Step())
	}

}
