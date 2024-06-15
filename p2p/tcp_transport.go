package p2p

import (
	"fmt"
	"net"
	"sync"
)

// It represents the remote node over a TCP established connection.
type TCPPeer struct{
	// conn is the underlying connection of peer
	conn net.Conn


	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		conn: 	conn,
		outbound: outbound,
	}
}

type TCPTransport struct{
	listenAddress string
	listener  net.Listener

	mu sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport{
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err!=nil{
		return err;
	}

	go t.startAcceptLoop()
	return nil;
}

func (t *TCPTransport) startAcceptLoop (){
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error : %s\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn){
	fmt.Printf("%+v\n", conn)
}