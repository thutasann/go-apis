package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thutasann/pokedexcli/internal/pokeapi"
)

// PokeDex Repl Configuration Type
type config struct {
	pokeapiClient       pokeapi.Client             // Poke API client
	nextLocationAreaURL *string                    // Next Location Area URL
	prevLocationAreaURL *string                    // Previous Location Area URL
	caughtPokemon       map[string]pokeapi.Pokemon // Catught Pokemon Map
}

// CLI Command Struct
type cliCommand struct {
	name        string                         // name of command
	description string                         // description of command
	callback    func(*config, ...string) error // callback of command
}

// PokeDex Repl Configurastion
var Config = config{
	pokeapiClient:       pokeapi.NewClient(time.Hour),
	nextLocationAreaURL: nil,
	prevLocationAreaURL: nil,
	caughtPokemon:       make(map[string]pokeapi.Pokemon),
}

// Start the Repl CLI
func StartRepl(cfg *config) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")

		scanner.Scan()
		text := scanner.Text()

		cleaned := CleanInput(text)
		if len(cleaned) == 0 {
			continue
		}

		commandName := cleaned[0]
		args := []string{}
		if len(cleaned) > 1 {
			args = cleaned[1:]
		}

		availableCommands := getCommands()
		command, ok := availableCommands[commandName]
		if !ok {
			fmt.Println("invalid command")
			continue
		}
		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Clean the Input
func CleanInput(str string) []string {
	lowered := strings.ToLower(str)
	words := strings.Fields(lowered)
	return words
}

// Get Available Commands
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Prints the help menu",
			callback:    CallbackHelp,
		},
		"map": {
			name:        "map",
			description: "List the next page of location areas",
			callback:    CallbackMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List the previous page of location areas",
			callback:    CallbackMapB,
		},
		"explore": {
			name:        "explore {location_area}",
			description: "List the pokemon in a location area",
			callback:    CallbackExplore,
		},
		"catch": {
			name:        "catch {name}",
			description: "Catch the pokemon with name",
			callback:    CallbackCatch,
		},
		"inspect": {
			name:        "inspect {name}",
			description: "inspect the pokemon with name",
			callback:    CallbackInspect,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Podedex",
			callback:    CallbackExit,
		},
	}
}
