package concurrencypatterns

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// PeerResponse represents a peer reply
type PeerResponse struct {
	Peer     string
	Response string
	Error    error
}

// callPeer preforms a TCP request with context
func callPeer(
	ctx context.Context,
	address string,
	message string,
	results chan<- PeerResponse,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	dialer := net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		results <- PeerResponse{Peer: address, Error: err}
		return
	}
	defer conn.Close()

	// set deadline from context (important)
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetDeadline(deadline)
	}

	// send request
	fmt.Fprintln(conn, message)

	// Read response
	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	if err != nil {
		results <- PeerResponse{Peer: address, Error: err}
		return
	}

	results <- PeerResponse{
		Peer:     address,
		Response: resp,
	}
}

func TCP_Context_Timeout() {
	peers := []string{
		"localhost:8080",
		"localhost:8081",
		"localhost:8082",
	}

	// Request-scoped timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := make(chan PeerResponse)
	var wg sync.WaitGroup

	for _, peer := range peers {
		wg.Add(1)
		go callPeer(ctx, peer, "ping", results, &wg)
	}

	// Close channel when all peers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Error != nil {
			fmt.Printf("Peer %s error: %v\n", res.Peer, res.Error)
			continue
		}
		fmt.Printf("Peer %s replied: %s", res.Peer, res.Response)
	}
}
