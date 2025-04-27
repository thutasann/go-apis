package main

import (
	"bytes"
	"encoding/binary"
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
	EncKey            []byte            // Encryption Key
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

	store  *Store        // File Server's Store
	quitch chan struct{} // Quit Channel
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
	Key string // Message File Key
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

// Get the file with key by Serve from the Local first
// If not found, fetching from the network
func (s *FileServer) Get(key string) (io.Reader, error) {
	if s.store.Has(key) {
		log.Printf("[Get] [%s] Serving file (%s) from local\n", s.Transport.Addr(), key)
		_, r, err := s.store.Read(key)
		return r, err
	}

	log.Printf("[Get] [%s] dont have file (%s) locally, fetching from network\n", s.Transport.Addr(), key)

	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 500)

	for _, peer := range s.peers {
		// First, read the file size so we can limit the amount of bytes that we read
		// from the connection, so it will not keep hanging
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)
		log.Println("[Get] fileSize --> ", fileSize)

		n, err := s.store.WriteDecrypt(s.EncKey, key, io.LimitReader(peer, fileSize))
		if err != nil {
			log.Printf("[Get] Store Write Error: %s\n", err)
			return nil, err
		}

		log.Printf("[Get] [%s] received (%d) bytes over the network from (%s): ", s.Transport.Addr(), n, peer.RemoteAddr())

		peer.CloseStream()
	}

	_, r, err := s.store.Read(key)
	return r, err
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
			Size: size + 16,
		},
	}

	if err := s.broadcast(&msg); err != nil {
		log.Println("[StoreData] boradcast error: ", err)
		return err
	}

	time.Sleep(time.Millisecond * 5)

	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	mw.Write([]byte{p2p.IncomingStream})
	n, err := CopyEncrypt(s.EncKey, fileBuffer, mw)
	if err != nil {
		return err
	}

	log.Printf("[Store] [%s] received and written bytes to disk: %d\n", s.Transport.Addr(), n)

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
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	return nil
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
		return s.handleMessageStoreFile(from, v)
	case MessageGetFile:
		return s.handleMessageGetFile(from, v)
	}
	return nil
}

// Handle Messag Get File from the `Get()` function
func (s *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	if !s.store.Has(msg.Key) {
		return fmt.Errorf("[handleMessageGetFile] [%s] need to serve file (%s) but it does not exist on disk", s.Transport.Addr(), msg.Key)
	}

	log.Printf("[handleMessageGetFile] serving file (%s) over the network\n", msg.Key)

	fileSize, r, err := s.store.Read(msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		log.Println("[handleMessageGetFile] closing read closer...")
		defer rc.Close()
	}

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer %s not in map", from)
	}

	// First send the "incomingStream" byte to the peer and then w
	// we can send the file size as int64
	peer.Send([]byte{p2p.IncomingStream})

	// Send the file Size
	binary.Write(peer, binary.LittleEndian, fileSize)

	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	log.Printf("[handleMessageGetFile] [%s] written (%d) bytes over the network to %s\n", s.Transport.Addr(), n, from)

	return nil
}

// Handle Message Store File and Write Data to Disk
func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	log.Printf("[handleMessageStoreFile] from: %+s, msg: %+v\n", from, msg)

	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) couldnot be found in the peer list", from)
	}

	size, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		log.Println("[handleMessageStoreFile] store write error: ", err)
		return err
	}

	log.Printf("[handleMessageStoreFile] [%s] written size: %v\n", s.Transport.Addr(), size)

	peer.CloseStream()

	return nil
}

// Initialize
func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
