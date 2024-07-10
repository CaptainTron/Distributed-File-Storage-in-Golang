package p2p

import "errors"

var ErrInvalidHandshake = errors.New("invalid handshake")
type HandshakeFunc func(Peer) error


func NOPHandshakeFunc(Peer) error { return nil }

