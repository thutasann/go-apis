package main

import (
	"fmt"
	"log"

	"github.com/thutasann/pokedexcli/internal/pokeapi"
)

// Pokedex CLI Tool
func main() {
	pokeapiClient := pokeapi.NewClient()
	resp, err := pokeapiClient.ListenLocationAreas()
	if err != nil {
		log.Fatal(resp)
	}
	fmt.Println(resp)
	// util.StartRepl()
}
