package main

import (
	"fmt"
	"log"

	"github.com/thutasann/distributed-file-storage/p2p"
)

func onPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

// Distributed File Storage
func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        onPeer,
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
