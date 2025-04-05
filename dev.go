package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Command struct {
	name    string
	command string
}

func main() {
	commands := []Command{
		{"Main Server", "npm start"},
		{"Socket Server", "npm run socket"},
		{"Auth Server", "npm run auth"},
		{"Server Manager", "npm run servers"},
		{"Channels Server", "npm run channels"},
		{"Payment Server", "npm run payment"},
		{"Notifications", "npm run notifications"},
		{"V2 Server", "npm run start:v2:watch"},
		{"V2 Socket", "npm run start:v2:socket:watch"},
		{"React App", "npm run react:dev"},
	}

	// Create new iTerm window
	exec.Command("osascript", "-e", `
		tell application "iTerm"
			create window with default profile
			tell current window
				tell current session
					split horizontally with default profile
					split vertically with default profile
				end tell
				tell second session
					split vertically with default profile
				end tell
			end tell
		end tell
	`).Run()

	time.Sleep(2 * time.Second)

	// Execute commands in each pane
	for i, cmd := range commands {
		applescript := fmt.Sprintf(`
			tell application "iTerm"
				tell current window
					tell session %d
						write text "cd %s"
						write text "%s"
					end tell
				end tell
			end tell
		`, i+1, os.Getenv("PWD"), cmd.command)

		exec.Command("osascript", "-e", applescript).Run()
		time.Sleep(1 * time.Second)
	}
}
