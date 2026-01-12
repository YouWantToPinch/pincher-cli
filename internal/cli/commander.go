package cli

import (
	"fmt"
	"log/slog"
	"strings"
)

// ================ PARSING USER INPUT ==================

// command represents a submission by a user through the CLI.
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
	optionsToParse := handler.options
	var parsingOption *cmdElement
	optArgCountNeeded := 0
	var actionElement *cmdElement

	allParamsSatisfied := func() bool {
		if actionElement != nil {
			return len(c.args) == handler.argCount()+actionElement.argCount()
		}
		return false
	}

	canTakeOpt := func() bool {
		return allParamsSatisfied() || len(c.args) == handler.argCount()
	}

	for i := 1; i < len(cmdFields); i++ {
		// are we parsing an option?
		if parsingOption != nil {
			// parsing an option; include in option's own argument stack
			c.opts[parsingOption.name] = append(c.opts[parsingOption.name], cmdFields[i])
			optArgCountNeeded -= 1
			if optArgCountNeeded == 0 {
				parsingOption = nil
			}
		} else {
			// have we encountered a potential option?
			if strings.HasPrefix(cmdFields[i], "--") ||
				(strings.HasPrefix(cmdFields[i], "-") && len(cmdFields[i]) == 2 && !strings.Contains("0123456789", string(cmdFields[i][1]))) {
				// can we take an option right now?
				if canTakeOpt() {
					// find out if the handler takes this option
					userOpt := strings.TrimLeft(cmdFields[i], "-")
					foundMatch := false
					for _, opt := range optionsToParse {
						foundMatch = (opt.name == userOpt || ((opt.letter() == userOpt) && opt.useShorthand))
						if foundMatch {
							c.opts[opt.name] = []string{}
							if opt.argCount() > 0 {
								parsingOption = &opt
								optArgCountNeeded = opt.argCount()
							} else {
								// if the option is valid, but takes no arguments, it must be a flag
								c.opts[opt.name] = append(c.opts[opt.name], "SET")
							}
							break
						}
					}
					// return error if this option is not taken by the handler
					if !foundMatch {
						optType := "command"
						if actionElement != nil {
							optType = "action"
						}
						return fmt.Errorf("input command includes unexpected %s option '%s'", optType, cmdFields[i])
					}
				}
			} else if allParamsSatisfied() {
				return fmt.Errorf("input command includes unexpected argument '%s'", cmdFields[i])
			}
			// not parsing an option; include in command argument stack
			c.args = append(c.args, cmdFields[i])
			// check whether or not we are past the point of parsing command options,
			// and onto the opportunity of parsing action options
			if el, found := findCMDElementWithName(handler.actions, cmdFields[i]); found {
				actionElement = el
				optionsToParse = el.options
			}
		}
	}
	if optArgCountNeeded > 0 {
		return fmt.Errorf("command could not be parsed; missing positional argument(s) for option [%s]: <%s>", parsingOption.name, parsingOption.parameters[len(c.opts[parsingOption.name])])
	}
	// if a command action was specified...
	if actionElement != nil {
		// check that the user has satisfied all parameters with arguments
		expectedArgs := append(handler.parameters, actionElement.parameters...)
		if len(c.args) < len(expectedArgs) {
			return fmt.Errorf("command could not be parsed; missing positional argument: <%s>", expectedArgs[len(c.args)])
		}
	}
	return nil
}

// ============== COMMAND ELEMENTS =================

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
	// A slice of expected arguments, in order, expected by this cmdElement.
	parameters []string
	// Options that this cmdElement accepts.
	// Should go UNUSED in cases where cmdElement IS an option, for obvious reasons.
	options []cmdElement
	// whether or not this element is an option that may be treated as a flag
	useShorthand bool
}

func (e *cmdElement) usage(withOptions bool) string {
	usage := e.name
	for _, arg := range e.parameters {
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
	return len(e.parameters)
}

func (e *cmdElement) letter() string {
	return string(e.name[0])
}

type HandlerFunc func(*State, *handlerContext) error

type handlerContext struct {
	cmd       command
	args      argTracker
	ctxValues map[string]string
}

// ========== ARGUMENT TRACKING =============

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

// cmdHandler represents a command which can be run.
type cmdHandler struct {
	cmdElement
	nonRegMsg string
	actions   []cmdElement
	callback  HandlerFunc
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

// ========== REGISTRY =============

type registrationStatus int

const (
	NotRegistered registrationStatus = iota
	Preregistered
	Registered
)

// a commandRegistry tracks any number of commands that a user
// has available to use. Some commands may require prerequisite
// actions before they are registered.
type commandRegistry struct {
	handlers map[string]*cmdHandler
	// preregistration of commands allows for users to predictably
	// know what commands are in theory POSSIBLE to execute within
	// the CLI, without yet registering them for use.
	// It means users can get output letting them know that the
	// command they entered isn't 'unknown', but that executing it
	// successfully is still dependent on some other action happening first,
	// such as logging in.
	registry map[string]registrationStatus
}

func (c *commandRegistry) run(s *State, input string) error {
	cmd := command{
		name: cleanInput(input)[0],
		opts: map[string][]string{},
	}

	// run the command only if it is fully registered
	status, registered := c.registry[cmd.name]
	if !registered {
		return fmt.Errorf("unknown command '%s'", cmd.name)
	}
	switch status {
	case Preregistered:
		return fmt.Errorf("cannot execute command '%s': %s", cmd.name, c.handlers[cmd.name].nonRegMsg)
	case NotRegistered:
		fallthrough
	default:
		return fmt.Errorf("unknown command '%s'", cmd.name)
	case Registered:
		// command is registered for use, so it may be run
	}

	handler, registered := c.handlers[cmd.name]
	if !registered {
		return fmt.Errorf("command found in registry, but without any handler")
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

// preregister establishes the existence of a command and its handler.
func (c *commandRegistry) preregister(handler *cmdHandler) {
	if handler == nil {
		slog.Error("denied preregistration of invalid handler")
		return
	}
	if handler.name == "" {
		slog.Error("denied preregistration of command without name")
		return
	}
	cmdName := handler.name

	_, registered := c.registry[cmdName]
	if !registered {
		c.handlers[cmdName] = handler
		c.registry[cmdName] = Preregistered
	} else {
		slog.Warn("attempted preregistration of command with handler, but is already preregistered", slog.String("command", cmdName))
	}
}

// register makes a preregistered command and its handlers available for use.
func (c *commandRegistry) register(cmdName string) {
	status, registered := c.registry[cmdName]
	if registered {
		switch status {
		case Preregistered:
			c.registry[cmdName] = Registered
		case Registered:
			slog.Warn("attempted registration of command for use, but is already registered", slog.String("command", cmdName))
		}
	}
}

// deregister looks for the command with the given name and,
// if it exists in the registry, updates its status to Preregistered.
func (c *commandRegistry) deregister(name string) {
	status, registered := c.registry[name]
	if registered {
		switch status {
		case Preregistered:
			slog.Warn("attempted deregistration of command for use, but is already in a preregistered state", slog.String("command", name))
		case Registered:
			c.registry[name] = Preregistered
			slog.Warn("attempted registration of command for use, but is already registered", slog.String("command", name))
		}
	}
}

func (c *commandRegistry) batchRegistration(handlers []*cmdHandler, newStatus registrationStatus) {
	for _, handler := range handlers {
		switch newStatus {
		case Preregistered:
			c.preregister(handler)
		case Registered:
			c.register(handler.name)
		}
	}
}

// deregisterNonBaseCommands deregisters all commands whose
// registration is dependent upon prior actions, such as
// logging in, or viewing a budget.
func (c *commandRegistry) deregisterNonBaseCommands() {
	for name, status := range c.registry {
		if status == Registered {
			c.deregister(name)
		}
	}
	c.batchRegistration(makeBaseCommandHandlers(), Registered)
}

func (c *commandRegistry) exists(name string) (*cmdHandler, bool) {
	handler, ok := c.handlers[name]
	if ok {
		return handler, ok
	}
	return nil, false
}

func (c *commandRegistry) GetRegisteredHandlers(verbose bool) []*cmdHandler {
	handlers := make([]*cmdHandler, 0, len(c.handlers))
	for cmdName, handler := range c.handlers {
		if verbose {
			handlers = append(handlers, handler)
		} else if status, exists := c.registry[cmdName]; exists && (status == Registered) {
			handlers = append(handlers, handler)
		}
	}
	return handlers
}
