package main

import (
	"fmt"
	"time"

	"github.com/thuta/go-elevator/internal/elevator"
)

func main() {
	e := elevator.New(0)

	e.AddRequest(3)
	e.AddRequest(1)
	e.AddRequest(5)

	tick := 0

	for {
		tick++

		fmt.Printf(
			"[tick %02d] floor=%d direction=%s targets=%v\n",
			tick,
			e.CurrentFloor(),
			e.Direction(),
			e.Targets(),
		)

		e.Step()

		if e.IsIdle() {
			fmt.Println("Elevator is idle. Simulation complete.")
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
}
