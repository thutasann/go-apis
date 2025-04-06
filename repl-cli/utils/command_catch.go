package util

import (
	"errors"
	"fmt"
	"math/rand"
)

// Callback Map to List Location Areas
func CallbackCatch(cfg *config, args ...string) error {

	if len(args) != 1 {
		return errors.New("no pokemon name provided")
	}

	pokemonName := args[0]

	pokemon, err := cfg.pokeapiClient.GetPokemon(&pokemonName)
	if err != nil {
		return err
	}

	const threshold = 50
	randNum := rand.Intn(pokemon.BaseExperience)
	if randNum < threshold {
		return fmt.Errorf("failed to catch %s", pokemonName)
	}

	fmt.Printf("Pokemon catched %s:\n", pokemon.Name)

	return nil
}
