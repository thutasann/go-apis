package util

import (
	"errors"
	"fmt"
)

// Callback Map to List Location Areas
func CallbackMap(cfg *config) error {
	resp, err := cfg.pokeapiClient.ListenLocationAreas(cfg.nextLocationAreaURL)
	if err != nil {
		return err
	}
	fmt.Println("Location area:")
	for _, area := range resp.Results {
		fmt.Printf(" - %s\n", area.Name)
	}
	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous
	return nil
}

// Callback Map to List Location Areas Back
func CallbackMapB(cfg *config) error {
	if cfg.prevLocationAreaURL != nil {
		return errors.New("you're on the first page")
	}

	resp, err := cfg.pokeapiClient.ListenLocationAreas(cfg.prevLocationAreaURL)
	if err != nil {
		return err
	}
	fmt.Println("Location area:")
	for _, area := range resp.Results {
		fmt.Printf(" - %s\n", area.Name)
	}
	cfg.nextLocationAreaURL = resp.Next
	cfg.prevLocationAreaURL = resp.Previous
	return nil
}
