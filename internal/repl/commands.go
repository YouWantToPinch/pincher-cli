package repl

import (
	"fmt"
)

// command represents a user-attempted input
type command struct {
	name string
	args []string
}

func (c *command) require(argCount int) error {
	if len(c.args) < argCount {
		return fmt.Errorf("Not enough arguments for command: %s", c.name)
	}
	return nil
}

type cmdFlag struct {
	word        string
	letter      string
	argCount    int64
	description string
	isOptional  bool
}

// cmdHandler represents a command which can be run.
type cmdHandler struct {
	name        string
	description string
	flags       []cmdFlag
	usage       string
	callback    func(*State, command) error
	// priority refers to a handler's relevance to users.
	// It is purely used for the purpose of sorting outputs related
	// to the handlers.
	// The lower the value, the higher the priority.
	//
	// Handlers more integral to the base functioning of the CLI,
	// such as 'exit', 'help', and 'config', reserve values 0-99.
	//
	// Handlers more relevant after a login, such as 'list',
	// reserve values over 99.
	//
	// The further away a command gets from relevance to the state of
	// the CLI at startup, the lower priority it ought to be given.
	priority int
}

func (c *cmdHandler) help() {
	fmt.Println("COMMAND: " + c.name)
	fmt.Println(c.description)
	if c.usage != "" {
		fmt.Println("USAGE: " + c.usage)
	}
	if len(c.flags) > 0 {
		fmt.Println("OPTIONS:")
		maxLen := MaxOfStrings(ExtractStrings(c.flags, func(c cmdFlag) string { return c.word }))
		for _, flag := range c.flags {
			fmt.Printf("  %-*s  %s\n", maxLen, flag.word, flag.description)
		}
	}
	fmt.Println()
}

type commandRegistry struct {
	handlers map[string]cmdHandler
}

func (c *commandRegistry) run(s *State, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Unknown command '%s'", cmd.name)
	}
	return handler.callback(s, cmd)
}

func (c *commandRegistry) register(name string, handler cmdHandler) {
	_, ok := c.handlers[name]
	if ok {
		fmt.Printf("ERROR: Command '%s' already exists in command registry\n", name)
	}
	c.handlers[name] = handler
}

func (c *commandRegistry) exists(name string) (cmdHandler, bool) {
	handler, ok := c.handlers[name]
	if ok {
		return handler, ok
	}
	return cmdHandler{}, false
}

func (c *commandRegistry) GetRegisteredHandlers() []cmdHandler {
	handlers := make([]cmdHandler, 0, len(c.handlers))
	for _, handler := range c.handlers {
		handlers = append(handlers, handler)
	}
	return handlers
}
