package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test TCP Transport
func TestTCPTransport(t *testing.T) {

	tcpOpts := TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: NOPHandShakeFunc,
		Decoder:       DefaultDecoder{},
	}

	listenAddr := ":3000"
	tr := NewTCPTransport(tcpOpts)
	assert.Equal(t, tr.ListenAddr, listenAddr)

	// server
	assert.Nil(t, tr.ListenAndAccept())

	select {}
}
