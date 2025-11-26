package repl

import (
	"fmt"
	"strings"
)

// command represents a user-attempted input
type command struct {
	name string
	// positional arguments that a handler may expect
	args []string
	// options, mapped to their own subset of arguments, that a handler may permit
	opts map[string][]string
}

func (c *command) hasOpt(name string) ([]string, bool) {
	opt, ok := c.opts[name]
	if ok {
		return opt, true
	}
	return []string{}, false
}

func (c *command) parse(h *cmdHandler, input string) error {

	for _, opt := range h.opts {
		if opt.isMandatory && !strings.Contains(input, "-"+opt.letter()) {
			return fmt.Errorf("ERROR: command could not be parsed; missing mandatory flag: '--%s'", opt.word)
		}
	}

	cmdFields := cleanInput(input)
	c.name = cmdFields[0]
	parsingFlag := ""
	argCountNeeded := 0

	for i := 1; i < len(cmdFields); i++ {
		// have we encountered a flag?
		if strings.HasPrefix(cmdFields[i], "-") {
			// return error if we encounter a flag while still parsing another one
			if parsingFlag != "" {
				return fmt.Errorf("ERROR: command could not be parsed; missing positional argument(s) for option: '--%s'", parsingFlag)
			}
			// find out of the handler takes this flag
			userOpt := strings.TrimLeft(cmdFields[i], "-")
			foundMatch := false
			for _, opt := range h.opts {
				fmt.Println("Checking: " + opt.word + " against " + userOpt)
				foundMatch = (opt.word == userOpt || opt.letter() == userOpt)
				if foundMatch {
					c.opts[opt.word] = []string{}
					if opt.argCount > 0 {
						parsingFlag = opt.word
						argCountNeeded = opt.argCount
					}
					break
				}
			}
			// return error if this flag is not taken by the handler
			if !foundMatch {
				return fmt.Errorf("ERROR: command includes unexpected option '%s'", cmdFields[i])
			}
		} else {
			if parsingFlag == "" {
				// not parsing a flag; include in arguments
				c.args = append(c.args, cmdFields[i])
			} else {
				// parsing a flag; include in flag's own argument stack
				c.opts[parsingFlag] = append(c.opts[parsingFlag], cmdFields[i])
				argCountNeeded -= 1
				if argCountNeeded == 0 {
					parsingFlag = ""
				}
			}
		}
	}
	if argCountNeeded > 0 {
		return fmt.Errorf("ERROR: command could not be parsed; missing positional argument(s) for option: '--%s'", parsingFlag)
	}
	return nil
}

func (c *command) require(argCount int) error {
	if len(c.args) < argCount {
		return fmt.Errorf("Not enough arguments for command: %s", c.name)
	}
	return nil
}

type cmdOption struct {
	word        string
	argCount    int
	description string
	isMandatory bool
}

func (c *cmdOption) letter() string {
	return string(c.word[0])
}

// cmdHandler represents a command which can be run.
type cmdHandler struct {
	name        string
	description string
	opts        []cmdOption
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
	if len(c.opts) > 0 {
		fmt.Println("OPTIONS:")
		maxLen := MaxOfStrings(ExtractStrings(c.opts, func(c cmdOption) string { return c.word }))
		for _, flag := range c.opts {
			fmt.Printf("  %-*s  %s\n", maxLen, flag.word, flag.description)
		}
	}
	fmt.Println()
}

type commandRegistry struct {
	handlers map[string]cmdHandler
}

func (c *commandRegistry) run(s *State, input string) error {
	cmd := command{opts: map[string][]string{}}

	handler, ok := c.handlers[cleanInput(input)[0]]
	if !ok {
		return fmt.Errorf("Unknown command '%s'", cmd.name)
	}

	err := cmd.parse(&handler, input)
	if err != nil {
		return err
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
