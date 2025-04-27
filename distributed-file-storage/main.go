package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/thutasann/distributed-file-storage/p2p"
)

// make new file server
func makeServer(listenAddr string, nodes ...string) *FileServer {

	// TCP Transport Options
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	// TCP Transport
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	// File Server Options
	fileServerOpts := FileServerOpts{
		EncKey:            NewEncryptionKey(),
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	// Create new File Server
	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

// Distributed File Storage
func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":7000", ":3000")
	s3 := makeServer(":5000", ":3000", ":7000")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()

	time.Sleep(2 * time.Second)

	go s3.Start()
	time.Sleep(2 * time.Second)

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("picture_%d.png", i)
		data := bytes.NewReader([]byte("my big data file here!"))
		s3.Store(key, data)

		if err := s3.store.Delete(key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("[main] found file --> ", string(b))
	}
}
