// Package repl contains all logic pertaining to the Pincher CLI's Read-Execute-Print loop,
// including handlers and commands.
package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
)

type State struct {
	Config          *config.Config
	Client          *client.Client
	CommandRegistry *commandRegistry
}

func StartRepl(cliState *State) {
	if cliState == nil {
		panic("StartRepl: cliState is nil")
	}

	mdAct := middlewareValidateAction

	cmdRegistry := &commandRegistry{
		handlers: make(map[string]*cmdHandler),
	}
	cmdRegistry.register("exit", &cmdHandler{
		cmdElement: cmdElement{
			name:        "exit",
			description: "exit the program",
			priority:    0,
		},
		callback: handlerExit,
	})
	cmdRegistry.register("help", &cmdHandler{
		cmdElement: cmdElement{
			name:        "help",
			description: "See usage of another command",
			arguments:   []string{"command"},
			priority:    1,
		},
		callback: handlerHelp,
	})
	cmdRegistry.register("config", &cmdHandler{
		cmdElement: cmdElement{
			name:        "config",
			description: "Add, Load, or Save a local user configuration for the Pincher-CLI",
			arguments:   []string{"action"},
			priority:    2,
		},
		actions: []cmdElement{
			{
				name:        "edit",
				description: "edit current user configuration",
			},
			{
				name:        "load",
				description: "load user configuration from the local machine",
			},
		},
		callback: mdAct(handlerConfig),
	})
	cmdRegistry.register("log", &cmdHandler{
		cmdElement: cmdElement{
			name:        "log",
			description: "see Pincher-CLI logs",
			priority:    3,
		},
		callback: handlerLog,
	})
	cmdRegistry.register("connect", &cmdHandler{
		cmdElement: cmdElement{
			name:        "connect",
			description: "Connect to a remote or local database",
			priority:    4,
		},
		callback: handlerConnect,
	})
	cmdRegistry.register("ready", &cmdHandler{
		cmdElement: cmdElement{
			name:        "ready",
			description: "Check server readiness",
			priority:    5,
		},
		callback: handlerReady,
	})

	cmdRegistry.register("user", &cmdHandler{
		cmdElement: cmdElement{
			name:        "user",
			description: "Create a new user, or log in",
			arguments:   []string{"action"},
			priority:    50,
		},
		callback: mdAct(handlerUser),
		actions: []cmdElement{
			{
				name:      "add",
				arguments: []string{"new_username", "new_password", "retype password"},
			},
			{
				name:      "login",
				arguments: []string{"username", "password"},
			},
		},
	})

	cliState.CommandRegistry = cmdRegistry
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
				fmt.Println(err.Error())
			}
		}
	}
}
