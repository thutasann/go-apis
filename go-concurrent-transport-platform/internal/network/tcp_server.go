package network

import "net"

// TCP Server wraps a TCP Listener
type TCPServer struct {
	Addr string
}

func (s *TCPServer) Listen() (net.Listener, error) {
	return net.Listen("tcp", s.Addr)
}
