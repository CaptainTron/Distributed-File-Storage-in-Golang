package main

import (
	"example/learn/p2p"
	"fmt"
	"log"
)

func OnPeer(p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside of TCP Transport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v/n", msg)
		}
	}()

	fmt.Println("Listening to New Request!!")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
