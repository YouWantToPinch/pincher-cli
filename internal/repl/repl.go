package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartRepl(cliState *State) {
	if cliState == nil {
		panic("startRepl: cliState is nil")
	}

	cmdRegistry := commandRegistry{
		handlers: make(map[string]cmdHandler),
	}
	//cmdRegistry.register("reset", handlerReset)
	cmdRegistry.register("exit", cmdHandler{
		name:        "exit",
		description: "exit the program",
		callback:    handlerExit,
	})
	cmdRegistry.register("help", cmdHandler{
		name:        "help",
		description: "See usage of the program",
		minArgs:     0,
		callback:    handlerHelp,
	})
	cmdRegistry.register("config", cmdHandler{
		name:        "config",
		description: "Add, Load, or Save a local user configuration for the Pincher-CLI",
		flags: []cmdFlag{
			cmdFlag{word: "edit", letter: "e",
				description: "edit current user configuration"},
			cmdFlag{word: "load", letter: "l",
				description: "load a user configuration from the local machine"},
		},
		minArgs:  1,
		callback: handlerConfig,
	})
	cmdRegistry.register("connect", cmdHandler{
		name:        "connect",
		description: "Conenct to a remote or local database",
		callback:    handlerConnect,
	})
	cmdRegistry.register("log", cmdHandler{
		name:        "log",
		description: "see Pincher-CLI logs",
		callback:    handlerLog,
	})
	cmdRegistry.register("report", cmdHandler{
		name:        "report",
		description: "Get a report from the database",
		callback:    handlerReport,
	})
	cmdRegistry.register("ready", cmdHandler{
		name:        "ready",
		description: "Check server readiness",
		callback:    handlerReady,
	})
	cmdRegistry.register("add", cmdHandler{
		name:        "add",
		description: "add an instance of some resource to the database",
		callback:    handlerAdd,
	})
	cmdRegistry.register("list", cmdHandler{
		name:        "list",
		description: "list instances of some resource from the database",
		callback:    handlerList,
	})

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the Pincher CLI!")
	fmt.Println("Use 'help' for available commands.")
	for {
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
