package repl

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
)

type State struct {
	Config *config.Config
	Client *client.Client
}

type command struct {
	name        string
	args        []string
	description string
}

type cmdFlag struct {
	word        string
	letter      string
	argCount    int64
	description string
	isOptional  bool
}

type cmdHandler struct {
	name        string
	flags       []cmdFlag
	description string
	callback    func(*State, command) error
	minArgs     int64
}

func (c *cmdHandler) parseFlags(args []string) {
	return
}

type commandRegistry struct {
	handlers map[string]cmdHandler
}

func (c *commandRegistry) run(s *State, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Unknown command '%s'", cmd.name)
	}
	if len(cmd.args) < int(handler.minArgs) {
		return fmt.Errorf("Not enough arguments for command: %s", handler.name)
	}
	return handler.callback(s, cmd)
}

func (c *commandRegistry) register(name string, handler cmdHandler) {
	_, ok := c.handlers[name]
	if ok {
		fmt.Println(fmt.Sprintf("ERROR: Command '%s' already exists in command registry", name))
	}
	c.handlers[name] = handler
	return
}

func (c *commandRegistry) exists(name string) error {
	_, ok := c.handlers[name]
	if ok {
		return nil
	}
	return fmt.Errorf("Command does not exist")
}
