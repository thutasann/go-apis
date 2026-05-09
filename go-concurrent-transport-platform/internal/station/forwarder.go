package station

import (
	"context"
	"log"

	"github.com/thutasann/ctp/internal/network"
	"github.com/thutasann/ctp/internal/telemetry"
)

// Forwarder pushes telemetry to control center
type Forwarder struct {
	Pool   *network.ConnPool
	Logger *log.Logger
}

func (f *Forwarder) Run(ctx context.Context, in <-chan telemetry.Telemetry) {
	for {
		select {
		case t := <-in:
			conn, err := f.Pool.Get(ctx)
			if err != nil {
				f.Logger.Println("dial error: ", err)
				continue
			}

			err = telemetry.Encode(conn, t)
			f.Pool.Put(conn)

			if err != nil {
				f.Logger.Println("send error: ", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
