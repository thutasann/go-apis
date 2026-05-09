package p2p

import "errors"

// ErrInvalidHandShake is returned if the handshake between
// the local and remote node could not be established
var ErrInvalidHandShake = errors.New("invalid handshake")

// HandShaker Function Signature that
type HandshakeFunc func(Peer) error

// No Operation HandShake Funcrtion
func NOPHandShakeFunc(Peer) error { return nil }
