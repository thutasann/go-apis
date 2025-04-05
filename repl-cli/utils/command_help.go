package util

import "fmt"

// Callback Help to List Helps Commands
func CallbackHelp() error {
	fmt.Println("Welcome to the Podedex help menu!")
	fmt.Println("Here are your available commands:")

	availableCommands := getCommands()
	for _, cmd := range availableCommands {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println("")
	return nil
}
