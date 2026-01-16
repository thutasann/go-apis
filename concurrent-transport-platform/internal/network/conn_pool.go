package network

import (
	"context"
	"net"
)

// ConnPool is a bounded TCP connection pool
type ConnPool struct {
	addr string
	pool chan net.Conn
}

// Initialize a new Connection Pool
func NewConnPool(addr string, size int) *ConnPool {
	return &ConnPool{
		addr: addr,
		pool: make(chan net.Conn, size),
	}
}

func (p *ConnPool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case c := <-p.pool:
		return c, nil
	default:
		var d net.Dialer
		return d.DialContext(ctx, "tcp", p.addr)
	}
}

func (p *ConnPool) Put(c net.Conn) {
	select {
	case p.pool <- c:
	default:
		c.Close()
	}
}
