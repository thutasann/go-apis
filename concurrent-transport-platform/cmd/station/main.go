package main

import (
	"context"

	"github.com/thutasann/ctp/internal/network"
	"github.com/thutasann/ctp/internal/shutdown"
	"github.com/thutasann/ctp/internal/station"
	"github.com/thutasann/ctp/internal/telemetry"
	"github.com/thutasann/ctp/pkg/logx"
)

func main() {
	ctx := shutdown.WithSignal(context.Background())

	logger := logx.New("STATION")

	buffer := station.NewBuffer[telemetry.Telemetry](1024)

	pool := network.NewConnPool("localhost:9000", 4)

	forwarder := &station.Forwarder{
		Pool:   pool,
		Logger: logger,
	}

	go forwarder.Run(ctx, buffer.Dequeue())

	srv := &station.Server{
		Addr:   ":8000",
		Buffer: buffer,
		Logger: logger,
	}

	_ = srv.Listen(ctx)

}
