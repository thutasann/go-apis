package concurrencypatterns

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// PooledConn wraps a net.Conn
type PooledConn struct {
	net.Conn
	inUse bool
}

// ConnPool manages reusable TCP connections
type ConnPool struct {
	address string
	pool    chan net.Conn
	dialer  net.Dialer
}

// Initialize a new Pool
func NewConnPool(address string, maxConns int) *ConnPool {
	return &ConnPool{
		address: address,
		pool:    make(chan net.Conn, maxConns),
		dialer:  net.Dialer{},
	}
}

// Get retrieves or creates a connection
func (p *ConnPool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case conn := <-p.pool:
		return conn, nil
	default:
		return p.dialer.DialContext(ctx, "tcp", p.address)
	}
}

// Put returns a connection to the pool
func (p *ConnPool) Put(conn net.Conn) {
	select {
	case p.pool <- conn:
	default:
		// Pool full -> close
		conn.Close()
	}
}

func Persistent_TCP_Peer_Connections() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := NewConnPool("localhost:8080", 2)

	var wg sync.WaitGroup

	sendRequest := func(id int) {
		defer wg.Done()

		conn, err := pool.Get(ctx)
		if err != nil {
			fmt.Println("Dial error:", err)
			return
		}

		// Impt: return connection to pool
		defer pool.Put(conn)

		conn.SetDeadline(time.Now().Add(2 * time.Second))

		fmt.Fprintf(conn, "request-%d\n", id)

		reader := bufio.NewReader(conn)
		resp, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		fmt.Printf("Response to %d: %s", id, resp)
	}

	// simulate concurrent requests
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go sendRequest(i)
	}

	wg.Wait()
}
