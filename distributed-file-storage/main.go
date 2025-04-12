package main

import (
	"fmt"
	"log"
	"time"

	"github.com/thutasann/distributed-file-storage/p2p"
)

// Distributed File Storage
func main() {
	fmt.Println("::: Starting Distributed File Storage :::")

	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":4000",
		HandshakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// todo: onPeer func
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       "4000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
