package util

import "os"

// Callback Exit to Exit the program
func CallbackExit(cfg *config) error {
	os.Exit(0)
	return nil
}
