package util

import "os"

// Callback Exit to Exit the program
func CallbackExit() error {
	os.Exit(0)
	return nil
}
