package train

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/thutasann/ctp/internal/network"
	"github.com/thutasann/ctp/internal/telemetry"
)

func RunTrain(
	ctx context.Context,
	addr string,
	trainID string,
	station string,
	logger *log.Logger,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := network.Dial(ctx, addr)
		if err != nil {
			logger.Println(trainID, "dail failed")
			time.Sleep(time.Second)
			continue
		}

		// simulate session
		sessionDuration := time.Duration(rand.Intn(10)+5) * time.Second
		end := time.Now().Add(sessionDuration)

		for time.Now().Before(end) {
			t := Generate(trainID, station)
			_ = telemetry.Encode(conn, t)

			time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
		}

		conn.Close()
		// simulate disconnect
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
	}
}
