package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WithSignal returns a context cancelled on SIGINT/SIGTERM
func WithSignal(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
