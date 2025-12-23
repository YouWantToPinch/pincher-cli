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

func (c *command) parse(handler *cmdHandler, input string) error {
	cmdFields := cleanInput(input)
	c.name = cmdFields[0]
	var parsingOptions []cmdElement
	parsingOptions = handler.options
	parsingOption := ""
	argCountNeeded := 0

	for i := 1; i < len(cmdFields); i++ {
		// have we encountered an option?
		if strings.HasPrefix(cmdFields[i], "-") {
			// return error if we encounter an option while still parsing another one
			if parsingOption != "" {
				return fmt.Errorf("ERROR: command could not be parsed; missing positional argument(s) for option: '--%s'", parsingOption)
			}
			// find out if the handler takes this option
			userOpt := strings.TrimLeft(cmdFields[i], "-")
			foundMatch := false
			for _, opt := range parsingOptions {
				// fmt.Println("Checking: " + opt.word + " against " + userOpt) // DEBUG
				foundMatch = (opt.name == userOpt || opt.letter() == userOpt)
				if foundMatch {
					c.opts[opt.name] = []string{}
					if opt.argCount() > 0 {
						parsingOption = opt.name
						argCountNeeded = opt.argCount()
					}
					break
				}
			}
			// return error if this option is not taken by the handler
			if !foundMatch {
				return fmt.Errorf("ERROR: command includes unexpected option '%s'", cmdFields[i])
			}
		} else {
			if parsingOption == "" {
				// not parsing an option; include in arguments
				c.args = append(c.args, cmdFields[i])
				// check whether or not we are past the point of parsing command options,
				// and onto the opportunity of parsing action options
				if el, found := findCMDElementWithName(handler.actions, cmdFields[i]); found {
					parsingOptions = el.options
				}
			} else {
				// parsing a option; include in option's own argument stack
				c.opts[parsingOption] = append(c.opts[parsingOption], cmdFields[i])
				argCountNeeded -= 1
				if argCountNeeded == 0 {
					parsingOption = ""
				}
			}
		}
	}
	if argCountNeeded > 0 {
		return fmt.Errorf("ERROR: command could not be parsed; missing positional argument(s) for option: '--%s'", parsingOption)
	}
	return nil
}

// cmdElement is the building block of the command infrastructure.
// It may represent any of:
//
// meta info for a command handler
//   - options accepted by a command handler
//   - actions accepted by a command handler
//     -> options accepted by an action
type cmdElement struct {
	name        string
	description string
	// Priority refers to an element's relevance to output.
	// The lower the value, the higher the priority.
	priority int
	// A slice of arguments, in order, expected by this cmdElement.
	arguments []string
	// Options that this cmdElement accepts.
	// Should go UNUSED in cases where cmdElement IS an option, for obvious reasons.
	options []cmdElement
}

func (e *cmdElement) usage(withOptions bool) string {
	usage := e.name
	for _, arg := range e.arguments {
		usage += " <" + arg + ">"
	}
	if len(e.options) > 0 {
		usage += " [options]\n"
		if withOptions {
			usage += "OPTIONS:\n"
			maxLen := MaxOfStrings(ExtractStrings(e.options, func(c cmdElement) string { return c.name }))
			for _, opt := range e.options {
				usage += fmt.Sprintf("  %-*s  %s\n", maxLen, opt.name, opt.description)
			}
		}
	} else {
		usage += "\n"
	}
	return usage
}

func (e *cmdElement) argCount() int {
	return len(e.arguments)
}

func (e *cmdElement) letter() string {
	return string(e.name[0])
}

type HandlerFunc func(*State, *handlerContext) error

// cmdHandler represents a command which can be run.
type cmdHandler struct {
	cmdElement
	actions  []cmdElement
	callback HandlerFunc
	//
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

type handlerContext struct {
	cmd       command
	args      argTracker
	ctxValues map[string]string
}

// argTracker is an iterator that progressively provides handlers with
// command argument OR option argument values
// (depending on the tracked index), with each call of its postfix function.
type argTracker struct {
	cmdArgs     *[]string
	cmdArgIndex int
	optArgs     *[]string
	optArgIndex int
}

func (a *argTracker) init(cmd *command) {
	a.reset()
	a.cmdArgs = &cmd.args
}

// start tracking the arguments of a given option
func (a *argTracker) trackOptArgs(cmd *command, option string) {
	optArgsAddressable := cmd.opts[option]
	a.optArgs = &optArgsAddressable
	a.optArgIndex = 0
}

// reset internal indeces
func (a *argTracker) reset() {
	a.cmdArgs = nil
	a.cmdArgIndex = 0
	a.optArgs = nil
	a.optArgIndex = 0
}

// return the currently-tracked index
func (a *argTracker) chooseIndex() (*int, *[]string) {
	if a.optArgs != nil {
		return &a.optArgIndex, a.optArgs
	}
	return &a.cmdArgIndex, a.cmdArgs
}

// postfix the tracked index, returning current value
func (a *argTracker) pfx() (string, error) {
	index, args := a.chooseIndex()
	if *index >= len(*args) {
		return "", fmt.Errorf("index out of range")
		// return fmt.Sprintf("ERROR: %d >= %d", *index, len(*args))
	}
	current := *index
	*index += 1
	return (*args)[current], nil
}

func (c *cmdHandler) help() {
	fmt.Println("COMMAND: " + c.name)
	fmt.Println(c.description)
	fmt.Println("USAGE: " + c.usage(true))
	if len(c.actions) > 0 {
		fmt.Println("ACTIONS:")
		fmt.Printf("(for further help, specify \"help %s <action>\")\n", c.name)
		maxLen := MaxOfStrings(ExtractStrings(c.actions, func(c cmdElement) string { return c.name }))
		for _, opt := range c.actions {
			fmt.Printf("  %-*s  %s\n", maxLen, opt.name, opt.description)
		}
	}
	fmt.Println()
}

type commandRegistry struct {
	handlers map[string]*cmdHandler
}

func (c *commandRegistry) run(s *State, input string) error {
	cmd := command{
		name: cleanInput(input)[0],
		opts: map[string][]string{},
	}

	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command '%s'", cmd.name)
	}

	err := cmd.parse(handler, input)
	if err != nil {
		return err
	}

	context := &handlerContext{
		cmd:       cmd,
		args:      argTracker{},
		ctxValues: make(map[string]string),
	}
	context.args.init(&context.cmd)

	return handler.callback(s, context)
}

func (c *commandRegistry) register(name string, handler *cmdHandler) {
	_, ok := c.handlers[name]
	if ok {
		fmt.Printf("ERROR: Command '%s' already exists in command registry\n", name)
	}
	c.handlers[name] = handler
}

func (c *commandRegistry) exists(name string) (*cmdHandler, bool) {
	handler, ok := c.handlers[name]
	if ok {
		return handler, ok
	}
	return nil, false
}

func (c *commandRegistry) GetRegisteredHandlers() []cmdHandler {
	handlers := make([]cmdHandler, 0, len(c.handlers))
	for _, handler := range c.handlers {
		handlers = append(handlers, *handler)
	}
	return handlers
}
