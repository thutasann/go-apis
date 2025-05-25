package main

import "fmt"

// Command Struct
type Command struct {
	//??
}

// Parse Command
func parseCommand(msg string) (Command, error) {
	t := msg[0]
	switch t {
	case '*':
		fmt.Println(msg[1:])
	}
	return Command{}, nil
}
