package p2p

import "net"

// it represents the connection b/w nodes
type Peer interface {
	net.Conn
	Send([]byte) error
}

// it will handle the connection between any nodes, like TCP, UDP and websockets
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
