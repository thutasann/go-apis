package main

import "github.com/thutasann/distributed-file-storage/p2p"

// File Server Options
type FileServerOpts struct {
	StorageRoot       string            // Storage Root
	PathTransformFunc PathTransformFunc // Path Transform Function
	Transport         p2p.Transport     // P2P Transport
}

// File Server Struct
type FileServer struct {
	FileServerOpts        // File Server Options
	store          *Store // File Server's Store
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
	}
}

// Start the File Server
func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}
