// Package cli contains all logic pertaining to the Pincher CLI state
package cli

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

	cliState.NewSession()

	// simulate login if session was saved
	if cliState.Config.StayLoggedIn && cliState.Client.RefreshToken != "" {
		cliState.Session.OnLogin()
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
			err := cliState.Session.CommandRegistry.run(cliState, input)
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
