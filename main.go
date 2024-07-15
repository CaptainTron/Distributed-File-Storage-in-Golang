package main

import (
	"bytes"
	"example/learn/p2p"
	// "fmt"
	// "io/ioutil"
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

	key := "myprivatedata"

	data := bytes.NewReader([]byte("Now, Stream of file is successfull and we can stream very large file"))
	s2.Store(key, data)

	// r, err := s2.Get(key)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// b, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(b))
	select {}
}

