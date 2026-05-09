package main

import (
	"context"

	"github.com/thutasann/ctp/internal/control"
	"github.com/thutasann/ctp/internal/shutdown"
)

func main() {
	ctx := shutdown.WithSignal(context.Background())

	agg := control.NewAggregator()

	srv := &control.Server{
		Addr:       ":9000",
		Aggregator: agg,
	}

	_ = srv.Listen(ctx)
}
