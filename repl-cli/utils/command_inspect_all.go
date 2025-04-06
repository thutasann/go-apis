package util

import (
	"fmt"
)

// Callback Map to Inspect Pokedex
func CallbackPokedex(cfg *config, args ...string) error {

	fmt.Println("Pokemon in Podedex:")
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}

	return nil
}
