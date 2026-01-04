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
		handlers: make(map[string]*cmdHandler),
		registry: make(map[string]registrationStatus),
	}
	cliState.CommandRegistry = cmdRegistry

	// preregister ALL commands
	cmdRegistry.batchRegistration(makeBaseCommandHandlers(), Preregistered)
	cmdRegistry.preregister(makeBudgetCommandHandler())
	cmdRegistry.batchRegistration(makeResourceCommandHandlers(), Preregistered)

	// register base commands

	cmdRegistry.batchRegistration(makeBaseCommandHandlers(), Registered)

	// fully register commands that require login if a session still exists
	if cliState.Config.StayLoggedIn && cliState.Client.LoggedInUser.RefreshToken != "" {
		cmdRegistry.register("budget")
	}

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
