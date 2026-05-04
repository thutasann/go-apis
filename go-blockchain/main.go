package main

import (
	"time"

	"github.com/thutasann/projectx/network"
)

func main() {
	trlocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trlocal.Connect(trRemote)
	trRemote.Connect(trlocal)

	go func() {
		for {
			trRemote.SendMessage(trlocal.Addr(), []byte("hello world"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trlocal},
	}

	s := network.NewServer(opts)
	s.Start()
}
