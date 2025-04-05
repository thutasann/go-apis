package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CLI Command Struct
type cliCommand struct {
	name        string       // name of command
	description string       // description of command
	callback    func() error // callback of command
}

// Start the Repl CLI
func StartRepl() {
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
		availableCommands := getCommands()
		command, ok := availableCommands[commandName]
		if !ok {
			fmt.Println("invalid command")
			continue
		}
		command.callback()
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
			description: "List some locatoin areas",
			callback:    CallbackMap,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Podedex",
			callback:    CallbackExit,
		},
	}
}
