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

// Message Struct
type Message struct {
	Payload any // Message Payload
}

// Message Store File
type MessageStoreFile struct {
	Key  string // Store File Key
	Size int64  // Byte size
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

// # Store File
//
// 1. Store the File to Disk
//
// 2. Broadcast this file to all known peers in the network
func (s *FileServer) StoreData(key string, r io.Reader) error {

	if err := s.store.Write(key, r); err != nil {
		fmt.Println("[StoreData] write error:", err)
		return err
	}

	buf := new(bytes.Buffer)
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: 22,
		},
	}

	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		fmt.Println("[StoreData] Encode Error: ", err)
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	for _, peer := range s.peers {
		n, err := io.Copy(peer, r)
		if err != nil {
			fmt.Printf("[StoreData] io copy error: %s\n", err)
			return err
		}
		fmt.Printf("[StoreData] received and written bytes to disk: %d\n", n)
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

// Broadcast the stored file to all known peers in the network
func (s *FileServer) Broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
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
				return
			}

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				fmt.Println("[loop] handleMessage Error : ", err)
				return
			}

		case <-s.quitch:
			return
		}
	}
}

// Handle Payload Message
func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleMessaegStoreFile(from, v)
	}
	return nil
}

// Handle Message Store File and Write Data to Disk
func (s *FileServer) handleMessaegStoreFile(from string, msg MessageStoreFile) error {
	log.Printf("[handleMessaegStoreFile] from: %+s, msg: %+v\n", from, msg)

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) couldnot be found in the peer list", from)
	}

	if err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		log.Println("[handleMessaegStoreFile] store write error: ", err)
		return err
	}

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

// Initialize
func init() {
	gob.Register(MessageStoreFile{})
}
