package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test TCP Transport
func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransport(listenAddr)
	assert.Equal(t, tr.listenAddress, listenAddr)

	// server
	assert.Nil(t, tr.ListenAndAccept())

	select {}
}
