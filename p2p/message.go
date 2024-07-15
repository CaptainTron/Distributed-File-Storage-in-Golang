package p2p

// import "net"
const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

//message represents data sent over network b/w nodes
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
