package mqtt

import (
	"context"
	"testing"

	"github.com/hemtjanst/bibliotek/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDiscover(t *testing.T) {
	if !enableIntegrationtests {
		t.Skip(integrationDisabledMsg)
	}

	hostPort, cleanup := testutils.MQTTBroker()
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	mq, err := New(ctx, &Config{Address: []string{hostPort}})
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	ch := mq.Discover()
	mq.Publish("discover", []byte("1"), true)
	v := <-ch
	assert.Equal(t, struct{}{}, v)
}

func TestResubscribe(t *testing.T) {
	if !enableIntegrationtests {
		t.Skip(integrationDisabledMsg)
	}

	hostPort, cleanup := testutils.MQTTBroker()
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	mq, err := New(ctx, &Config{Address: []string{hostPort}})
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	t.Run("without sub", func(t *testing.T) {
		ok := mq.Resubscribe("testing1", "testing2")
		assert.False(t, ok)
	})
	t.Run("with sub", func(t *testing.T) {
		_ = mq.Subscribe("testing1")
		ok := mq.Resubscribe("testing1", "testing2")
		assert.True(t, ok)
	})
}

func TestUnsubscribe(t *testing.T) {
	if !enableIntegrationtests {
		t.Skip(integrationDisabledMsg)
	}

	hostPort, cleanup := testutils.MQTTBroker()
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	mq, err := New(ctx, &Config{Address: []string{hostPort}})
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	t.Run("without sub", func(t *testing.T) {
		ok := mq.Unsubscribe("testing1")
		assert.False(t, ok)
	})
	t.Run("with sub", func(t *testing.T) {
		_ = mq.Subscribe("testing1")
		ok := mq.Unsubscribe("testing1")
		assert.True(t, ok)
	})
}
