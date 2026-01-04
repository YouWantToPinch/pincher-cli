// Package repl contains all logic pertaining to the Pincher CLI's Read-Execute-Print loop,
// including handlers and commands.
package repl

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

func StartRepl(cliState *State) {
	if cliState == nil {
		panic("StartRepl: cliState is nil")
	}

	err := cliState.Client.NewTokenOrLogout()
	if err != nil {
		slog.Info("logged out expired user session from cache")
	}

	cmdRegistry := &commandRegistry{
		handlers:      make(map[string]*cmdHandler),
		preregistered: make(map[string]bool),
	}
	cliState.CommandRegistry = cmdRegistry

	registerBaseCommands(cliState, false)

	// register commands that require login if a session still exists
	if cliState.Config.StayLoggedIn && cliState.Client.LoggedInUser.RefreshToken != "" {
		registerBudgetCommand(cliState, false)
	} else {
		registerBudgetCommand(cliState, true)
	}
	registerResourceCommands(cliState, true)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the Pincher CLI!")
	fmt.Println("Use 'help' for available commands.")
	for {
		fmt.Println("__________________")
		fmt.Print("Pincher > ")
		if scanner.Scan() {
			input := scanner.Text()
			if len(input) == 0 {
				continue
			}
			err := cmdRegistry.run(cliState, input)
			if err != nil {
				if err.Error() == "HIJACK:EXIT" {
					break
				}
				slog.Error(err.Error())
				fmt.Println("ERROR: " + err.Error())
			}
		}
	}
}
