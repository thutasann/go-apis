package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

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

// Broadcast Payload
type Payload struct {
	Key  string // Payload File Key
	Data []byte // Payload Data
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
func (s *FileServer) broadcast(p *Payload) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

// # Store File
//
// 1. Store the File to Disk
//
// 2. Broadcast this file to all known peers in the network
func (s *FileServer) StoreData(key string, r io.Reader) error {
	if err := s.store.Write(key, r); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}
	fmt.Println("[StoreData] byte --> ", buf.Bytes())

	p := &Payload{
		Key:  key,
		Data: buf.Bytes(),
	}

	return s.broadcast(p)
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
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}
