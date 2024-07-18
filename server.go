package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"example/learn/p2p"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type FileServerOpts struct {
	// ID will eventually help us to identify the nodes,
	// We can sync the whole folder, located at different nodes
	ID                string
	Encyption_key     []byte // for encryptionMethod
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	if len(opts.ID) == 0 {
		opts.ID = generateID()
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	ID   string
	Key  string
	Size int64
}

type MessageGetFile struct {
	ID  string
	Key string
}

// broadcast to all the available peers at once at the stream
func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	fmt.Printf("Starting BroadCast\n")
	for _, peer := range s.peers {
		// Send the message 
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(s.ID, key) {
		fmt.Printf("[%s] serving file (%s) from local disk\n", s.Transport.Addr(), key)
		_, r, err := s.store.Read(s.ID, key)
		return r, err
	}

	fmt.Printf("[%s] Don't have file (%s) locally, fetching from network\n", s.Transport.Addr(), key)
	msg := Message{
		Payload: MessageGetFile{
			ID:  s.ID,
			Key: hashkey(key),
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 5)
	// After receive from remote peer save into disk
	for _, peer := range s.peers {
		// First of all read the file size so we can limit the amount of bytes that we read
		// from connection, so it will not keep hanging.
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)
		// Storing into the disk after receiving from remote server...
		n, err := s.store.WriteDecrypt(s.Encyption_key, s.ID, key, io.LimitReader(peer, fileSize))
		//n, err := s.store.Write(key, io.LimitReader(peer, fileSize))
		if err != nil {
			return nil, err
		}
		fmt.Printf("[%s] Received %d bytes over the network from [%s]", s.Transport.Addr(), n, peer.RemoteAddr())
		peer.CloseStream()
	}
	_, r, err := s.store.Read(s.ID, key)
	return r, err
}

func (s *FileServer) Store(key string, r io.Reader) error {
	// 1. Store this file to disk
	// 2. Broadcast this file to all known peers in the network
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)

	// Save in this machine
	size, err := s.store.Write(s.ID, key, tee)
	if err != nil {
		return err
	}

	// Create message
	msg := Message{
		Payload: MessageStoreFile{
			ID:   s.ID,
			Key:  hashkey(key),
			Size: size + 16,
		},
	}

	// broadcast the message in every peer...
	if err := s.broadcast(&msg); err != nil {
		return err
	}

	fmt.Printf("Starting File upload\n")
	time.Sleep(time.Millisecond * 5)

	// Sending the file to all the peers at once

	// creating slice of writers, appending all the avaibable peers into it
	// creating multiwriter mw that will write into all the peers at once
	// at end send stream message, and send the encrypted message to all peers
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	mw.Write([]byte{p2p.IncomingStream})
	n, err := copyEncrypt(s.Encyption_key, fileBuffer, mw)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] send and written (%d) bytes to disk\n", s.Transport.Addr(), n)

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s", p.RemoteAddr().String())
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to error or user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			fmt.Println("Checking rpc:-------> on port ", s.Transport.Addr())
			var msg Message
			// Message will printed here
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("decoding err: ", err)
			}

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handleMessage Error: ", err)
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessageStoreFile(from, v)
	case MessageGetFile:
		return s.handleMessageGetFile(from, v)
	}
	return nil
}

// this will be executed in remote server call
func (s *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	if !s.store.Has(msg.ID, msg.Key) {
		return fmt.Errorf("[%s] need to serve file %s but it does not exists on disk", s.Transport.Addr(), msg.Key)
	}
	log.Printf("[%s] serving file (%s) over the network\n", s.Transport.Addr(), msg.Key)

	fileSize, r, err := s.store.Read(msg.ID, msg.Key)
	if err != nil {
		return err
	}

	// Check if it is ending or not
	// if reader is readcloser or not
	if rc, ok := r.(io.ReadCloser); ok {
		fmt.Println("Closing readcloser")
		defer rc.Close()
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("could not find the [%s] node in network", from)
	}

	// Send the incoming stream to the requesting server...
	peer.Send([]byte{p2p.IncomingStream})

	// Send the file to the remote server...
	binary.Write(peer, binary.NativeEndian, fileSize)
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] written %d bytes on the network to [%s]\n", s.Transport.Addr(), n, from)
	return nil
}

// StoreMessage [RemoteServer]
func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer: [%s] not found in peer list", from)
	}
	// Write in Remote Server, received from remote server...
	n, err := s.store.Write(msg.ID, msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}
	log.Printf("Remote Server.... written (%d) bytes to disk", n)
	peer.CloseStream()
	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		// Start a new Connection in another goroutine
		go func(addr string) {
			log.Printf("[%s] Attempting to connect with remote server %s\n", s.Transport.Addr(), addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)

	}
	return nil
}

func (s *FileServer) Start() error {
	log.Printf("Starting File Server %s", s.Transport.Addr())
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()
	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
