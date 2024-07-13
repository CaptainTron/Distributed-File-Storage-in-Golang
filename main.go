package main

import (
	"bytes"
	"example/learn/p2p"
	"log"
	"strings"
	"time"
	// "time"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO onPeerFunc
	}
	tcptransport := p2p.NewTCPTransport(tcptransportOpts)
	FileServerOpts := FileServerOpts{
		StorageRoot:       "dir_" + strings.TrimPrefix(listenAddr, ":"),
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcptransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(FileServerOpts)
	tcptransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(1 * time.Second)
	go s2.Start()
	time.Sleep(1 * time.Second)

	data := bytes.NewReader([]byte("my big data file here!"))
	s2.StoreData("myprivatedata", data)
	select {}
}
