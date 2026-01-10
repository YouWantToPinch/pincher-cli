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

	commandQueue := make(chan string, 32)
	cliState.CmdQueue = commandQueue

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the REPL CLI!")
	fmt.Println("Use 'help' for available commands.")
	for {
		fmt.Println("__________________")
		fmt.Print("Pincher > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if len(input) == 0 {
			continue
		}
		cliState.CmdQueue <- input

		for len(commandQueue) > 0 {
			cmd := <-commandQueue
			if cmd == "exit" {
				fmt.Println("Exiting Pincher CLI Program...")
				*cliState.DoneChan <- true
				return
			}
			err := cliState.Session.CommandRegistry.run(cliState, cmd)
			if err != nil {
				slog.Error(err.Error())
				fmt.Println("ERROR:", err)
			}
		}
	}
}
