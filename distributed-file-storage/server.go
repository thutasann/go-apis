package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/thutasann/distributed-file-storage/p2p"
)

// File Server Options
type FileServerOpts struct {
	StorageRoot       string            // Storage Root
	PathTransformFunc PathTransformFunc // Path Transform Function
	Transport         p2p.Transport     // P2P Transport
	BootstrapNodes    []string          // Bootstrap Nodes Arrays
}

// File Server Struct
type FileServer struct {
	FileServerOpts // File Server Options

	peerLock sync.Mutex          // Peer Lock
	peers    map[string]p2p.Peer // Peers Map
	store    *Store              // File Server's Store
	quitch   chan struct{}       // Quit Channel
}

// Message
type Message struct {
	From    string
	Payload any
}

// Initialize New File Server
func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

// Start the File Server
func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil
}

// Broadcast the stored file to all known peers in the network
func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

// # Store File
//
// 1. Store the File to Disk
//
// 2. Broadcast this file to all known peers in the network
func (s *FileServer) StoreData(key string, r io.Reader) error {
	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }

	// return s.broadcast(&Message{
	// 	From:    "todo",
	// 	Payload: p,
	// })

	buf := new(bytes.Buffer)
	msg := Message{
		Payload: []byte("storagekey"),
	}
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	payload := []byte("THIS LARGE FILE!")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}

	return nil
}

// Stop the File Server
// Close the Quit Channel
func (s *FileServer) Stop() {
	close(s.quitch)
}

// OnPeer function for the file Server
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("[OnPeer] connected with remote %s", p.RemoteAddr())
	return nil
}

// Bootstrap Networks
func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			log.Printf("[bootstrapNetwork] attempting to connect with remote: %s", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error ::> ", err)
			}
		}(addr)
	}
	return nil
}

// Loop the incoming messages and Consume
func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("[loop] Decode error : ", err)
			}

			fmt.Printf("[loop] recv: %s\n", string(msg.Payload.([]byte)))

			peer, ok := s.peers[rpc.From]
			if !ok {
				panic("[loop] peer is not found in peers map")
			}
			fmt.Println("[loop] peer -> ", peer)

			b := make([]byte, 1000)
			if _, err := peer.Read(b); err != nil {
				panic(err)
			}

			fmt.Printf("[loop] after read byte: %s\n", string(b))

			// if err := s.handleMessage(&m); err != nil {
			// 	log.Println("[loop] handleMessage error : ", err)
			// }

		case <-s.quitch:
			return
		}
	}
}

// Handle Payload Message
// func (s *FileServer) handleMessage(msg *Message) error {
// 	switch v := msg.Payload.(type) {
// 	case *DataMessage:
// 		fmt.Printf("received data %+v\n", v)
// 	}
// 	return nil
// }
