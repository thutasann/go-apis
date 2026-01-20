package train

import (
	"math/rand"
	"time"

	"github.com/thutasann/ctp/internal/telemetry"
)

func Generate(trainID, station string) telemetry.Telemetry {
	return telemetry.Telemetry{
		TrainID:        trainID,
		Station:        station,
		SpeedKmh:       rand.Intn(80) + 20,
		PassengerCount: rand.Intn(1500),
		Timestamp:      time.Now(),
	}
}
