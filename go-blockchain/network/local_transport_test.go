package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	assert.NoError(t, tra.Connect(trb))
	assert.NoError(t, trb.Connect(tra))

	message := []byte("ping")
	assert.NoError(t, tra.SendMessage(trb.Addr(), message))

	rpc := <-trb.Consume()
	assert.Equal(t, message, rpc.Payload)
	assert.Equal(t, tra.Addr(), rpc.From)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	assert.NoError(t, tra.Connect(trb))
	assert.NoError(t, trb.Connect(tra))

	msg := []byte("hello world")
	assert.NoError(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()
	assert.Equal(t, msg, rpc.Payload)
	assert.Equal(t, tra.Addr(), rpc.From)
}
