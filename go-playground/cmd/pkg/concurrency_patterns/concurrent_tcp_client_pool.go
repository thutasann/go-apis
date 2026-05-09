package concurrencypatterns

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

// PeerResult holds response from a peer
type PeerResult struct {
	Address  string
	Response string
	Error    error
}

// handles one TCP Peer
func connectToPeer(
	address string,
	message string,
	timeout time.Duration,
	results chan<- PeerResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		results <- PeerResult{
			Address: address,
			Error:   err,
		}
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	// Send Message
	fmt.Fprintln(conn, message)

	// Read Response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		results <- PeerResult{
			Address: address,
			Error:   err,
		}
		return
	}

	results <- PeerResult{
		Address:  address,
		Response: response,
	}
}

func Concurrent_TCP_Client_Pool() {
	peers := []string{
		"localhost:8080",
		"localhost:8081",
		"localhost:8082",
	}

	message := "ping"
	timeout := 2 * time.Second

	results := make(chan PeerResult)
	var wg sync.WaitGroup

	for _, peer := range peers {
		wg.Add(1)
		go connectToPeer(peer, message, timeout, results, &wg)
	}

	// Close results when all peers handled
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect responses
	for res := range results {
		if res.Error != nil {
			fmt.Printf("Peer %s error: %v\n", res.Address, res.Error)
			continue
		}

		fmt.Printf("Peer %s replied: %s", res.Address, res.Response)
	}
}
