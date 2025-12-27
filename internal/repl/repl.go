// Package repl contains all logic pertaining to the Pincher CLI's Read-Execute-Print loop,
// including handlers and commands.
package repl

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
)

type State struct {
	DoneChan        *chan bool
	Logger          *Logger
	Config          *config.Config
	Client          *client.Client
	CommandRegistry *commandRegistry
}

func StartRepl(cliState *State) {
	if cliState == nil {
		panic("StartRepl: cliState is nil")
	}

	cmdRegistry := &commandRegistry{
		handlers:      make(map[string]*cmdHandler),
		preregistered: make(map[string]bool),
	}
	cliState.CommandRegistry = cmdRegistry

	registerBaseCommands(cliState, false)
	registerBudgetCommand(cliState, true)
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
