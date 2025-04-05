package util

import (
	"fmt"
	"log"

	"github.com/thutasann/pokedexcli/internal/pokeapi"
)

// Callback Map to List Location Areas
func CallbackMap() error {
	pokeapiClient := pokeapi.NewClient()

	resp, err := pokeapiClient.ListenLocationAreas()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Location area:")
	for _, area := range resp.Results {
		fmt.Printf(" - %s\n", area.Name)
	}
	return nil
}
