package http

import (
	"net"
	"net/http"
	"time"
)

// NewServer creates tuned HTTP server.
//
// Tuning matters for high concurrency.
func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,

		// Custom listener tuning handled outside if needed
	}
}

// TuneListener sets low-level socket options.
func TuneListener(addr string) (net.Listener, error) {
	lc := net.ListenConfig{
		KeepAlive: 30 * time.Second,
	}
	return lc.Listen(nil, "tcp", addr)
}
