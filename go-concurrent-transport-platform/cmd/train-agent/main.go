package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/thutasann/ctp/internal/shutdown"
	"github.com/thutasann/ctp/internal/train"
	"github.com/thutasann/ctp/pkg/logx"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx := shutdown.WithSignal(context.Background())

	logger := logx.New("TRAIN")

	stationAddr := "localhost:8000"
	stations := []string{"Bugis", "CityHall", "Raffles"}

	trainCount := 100

	for i := 0; i < trainCount; i++ {
		trainID := fmt.Sprintf("EW-%03d", i)
		station := stations[rand.Intn(len(stations))]

		go train.RunTrain(ctx, stationAddr, trainID, station, logger)
	}

	<-ctx.Done()
}
