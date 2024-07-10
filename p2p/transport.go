package p2p

// it represents the connection b/w nodes
type Peer interface {
	Close() error
}

// it will handle the connection between any nodes, like TCP, UDP and websockets
type Transport interface {
	ListenAndAccept() error
	Consume() <- chan RPC
}
