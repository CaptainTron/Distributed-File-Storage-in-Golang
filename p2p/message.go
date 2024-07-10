package p2p

import "net"

//message represents data sent over network b/w nodes
type RPC struct {
	Payload []byte
	From net.Addr
}
