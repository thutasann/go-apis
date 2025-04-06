package util

import (
	"errors"
	"fmt"
)

// Callback Map to Inspect Pokemon After Catched
func CallbackInspect(cfg *config, args ...string) error {

	if len(args) != 1 {
		return errors.New("no pokemon name provided")
	}

	pokemonName := args[0]

	caughtedPokemon, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		return errors.New("you haven't caught this pokemon yet")
	}

	fmt.Printf("Inspected Pokemon Name : %s", caughtedPokemon.Name)

	return nil
}
