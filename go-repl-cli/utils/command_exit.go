package util

import "os"

// Callback Exit to Exit the program
func CallbackExit(cfg *config, args ...string) error {
	os.Exit(0)
	return nil
}
