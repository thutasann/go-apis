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
	FileServerOpts                     // File Server Options
	peerLock       sync.Mutex          // Peer Lock
	peers          map[string]p2p.Peer // Peers Map
	store          *Store              // File Server's Store
	quitch         chan struct{}       // Quit Channel
}

// Message Struct
type Message struct {
	Payload any // Message Payload
}

// Message Store File Struct
type MessageStoreFile struct {
	Key  string // Store File Key
	Size int64  // Byte size
}

// Message Get File Struct
type MessageGetFile struct {
	Key string
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

// Get the file with key
func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		return s.store.Read(key)
	}

	log.Printf("[Get] dont have file (%s) locally, fetching from network\n", key)

	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	select {}
}

// Store the File to Disk and broadcast this file to all known peers in the network
func (s *FileServer) Store(key string, r io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)

	size, err := s.store.Write(key, tee)
	if err != nil {
		fmt.Println("[StoreData] write error:", err)
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		log.Println("[StoreData] boradcast error: ", err)
		return err
	}

	time.Sleep(time.Second * 3)

	for _, peer := range s.peers {
		n, err := io.Copy(peer, fileBuffer)
		if err != nil {
			fmt.Printf("[StoreData] io copy error: %s\n", err)
			return err
		}
		log.Printf("[StoreData] received and written bytes to disk: %d\n", n)
	}

	return nil
}

// Stop the File Server and Close the Quit Channel
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

// broadcast the message
func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		fmt.Println("[StoreData] Encode Error: ", err)
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
}

// stream the stored file to all known peers in the network
func (s *FileServer) Stream(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

// Loop the incoming messages and Consume
func (s *FileServer) loop() {
	defer func() {
		log.Println("[loop] file server stopped due to error or user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("[loop] Decode error : ", err)
			}

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("[loop] handleMessage Error:", err)
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
	case MessageGetFile:
		return s.handleMessaegGetFile(from, v)
	}
	return nil
}

// Handle Messag Get File from the `Get()` function
func (s *FileServer) handleMessaegGetFile(from string, msg MessageGetFile) error {
	if !s.store.Has(msg.Key) {
		return fmt.Errorf("need to serve file (%s) but it does not exist on disk", msg.Key)
	}

	r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer %s not in map", from)
	}

	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	log.Printf("[handleMessaegGetFile] written %d bytes over the network to %s\n", n, from)

	return nil
}

// Handle Message Store File and Write Data to Disk
func (s *FileServer) handleMessaegStoreFile(from string, msg MessageStoreFile) error {
	log.Printf("[handleMessaegStoreFile] from: %+s, msg: %+v\n", from, msg)

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) couldnot be found in the peer list", from)
	}

	size, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		log.Println("[handleMessaegStoreFile] store write error: ", err)
		return err
	}

	log.Println("[handleMessaegStoreFile] written size: ", size)

	peer.(*p2p.TCPPeer).Wg.Done()

	return nil
}

// Initialize
func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
