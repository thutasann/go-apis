package main

import (
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
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(1 * time.Second)

	go s2.Start()
	time.Sleep(1 * time.Second)

	// data := bytes.NewReader([]byte("my big data file here!"))
	// s2.Store("coolPicture.jpg", data)
	// time.Sleep(5 * time.Millisecond)

	r, err := s2.Get("coolPicture.jpg")
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("found file --> ", string(b))
}
