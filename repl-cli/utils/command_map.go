package util

import (
	"fmt"
	"log"
)

// Callback Map to List Location Areas
func CallbackMap(cfg *config) error {
	resp, err := cfg.pokeapiClient.ListenLocationAreas(cfg.nextLocationAreaURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Location area:")
	for _, area := range resp.Results {
		fmt.Printf(" - %s\n", area.Name)
	}
	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous
	return nil
}
