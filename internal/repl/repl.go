package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		panic("startRepl: cliState is nil")
	}

	cmdRegistry := &commandRegistry{
		handlers: make(map[string]cmdHandler),
	}
	cmdRegistry.register("exit", cmdHandler{
		name:        "exit",
		description: "exit the program",
		priority:    0,
		callback:    handlerExit,
	})
	cmdRegistry.register("help", cmdHandler{
		name:        "help",
		description: "See usage of another command\n",
		usage:       "help <command>",
		priority:    1,
		callback:    handlerHelp,
	})
	cmdRegistry.register("config", cmdHandler{
		name:        "config",
		description: "Add, Load, or Save a local user configuration for the Pincher-CLI",
		usage:       "config (--edit | --load)",
		flags: []cmdFlag{
			{word: "edit", letter: "e",
				description: "edit current user configuration"},
			{word: "load", letter: "l",
				description: "load user configuration from the local machine"},
		},
		priority: 2,
		callback: handlerConfig,
	})
	cmdRegistry.register("log", cmdHandler{
		name:        "log",
		description: "see Pincher-CLI logs",
		priority:    3,
		callback:    handlerLog,
	})
	cmdRegistry.register("connect", cmdHandler{
		name:        "connect",
		description: "Conenct to a remote or local database",
		priority:    4,
		callback:    handlerConnect,
	})
	cmdRegistry.register("ready", cmdHandler{
		name:        "ready",
		description: "Check server readiness",
		priority:    5,
		callback:    handlerReady,
	})

	cmdRegistry.register("user", cmdHandler{
		name:        "user",
		description: "Create a new user, or log in",
		priority:    50,
		callback:    handlerUser,
		usage: `user <action> [options] [arguments]
user add <new_username> <new_password> <new_password>
user login <username> <password>`,
	})

	cmdRegistry.register("report", cmdHandler{
		name:        "report",
		description: "Get a report from the database",
		priority:    100,
		callback:    handlerReport,
	})
	cmdRegistry.register("list", cmdHandler{
		name:        "list",
		description: "list instances of some resource from the database",
		priority:    101,
		callback:    handlerList,
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
			command := command{name: cleanInput(input)[0], args: cleanInput(input)[1:]}
			err := cmdRegistry.run(cliState, command)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}
