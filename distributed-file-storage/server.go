package main

import (
	"fmt"
	"log"

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
	FileServerOpts               // File Server Options
	store          *Store        // File Server's Store
	quitch         chan struct{} // Quit Channel
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

// Stop the File Server
// Close the Quit Channel
func (s *FileServer) Stop() {
	close(s.quitch)
}

// Bootstrap Networks
func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		go func(addr string) {
			fmt.Println("attempting to connect with remote..", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
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
