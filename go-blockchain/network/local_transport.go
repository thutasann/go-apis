// Package network provides the transport layer for blockchain peer-to-peer communication.
// It defines the interfaces and implementations for peers to communicate with each other
// in a decentralized network. LocalTransport is an in-memory implementation suitable for
// testing and local development.
package network

import (
	"fmt"
	"sync"
)

// LocalTransport implements the Transport interface using in-memory channels.
// It maintains a map of connected peers and handles message passing between them using goroutines.
type LocalTransport struct {
	// addr is the unique network address of this transport peer.
	addr NetAddr
	// consumeCh is a buffered channel through which incoming RPC messages are delivered.
	consumeCh chan RPC
	// lock protects concurrent access to the peers map.
	lock sync.RWMutex
	// peers is a map of connected peers indexed by their network address.
	peers map[NetAddr]*LocalTransport
}

// NewLocalTransport creates and returns a new LocalTransport instance with the given address.
// The returned transport is ready to accept connections and send/receive messages.
func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
	}
}

// Consume returns a read-only channel for receiving incoming RPC messages.
// The caller should read from this channel in a separate goroutine.
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// Connect establishes a connection to another peer transport and adds it to the peers map.
// This allows this transport to send messages to the connected peer.
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

// Addr returns the network address of this transport peer.
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

// SendMessage sends a message payload to a peer at the specified address.
// It returns an error if the peer is not connected or not found.
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    t.addr,
		Payload: payload,
	}

	return nil
}
