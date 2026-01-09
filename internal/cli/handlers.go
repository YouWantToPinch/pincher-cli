package cli

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
)

// ========= MIDDLEWARE ==============

func middlewareValidateAction(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(s *State, c *handlerContext) error {
		var err error
		action, _ := c.args.pfx()
		if action == "" {
			err = fmt.Errorf("no action specified")
			s.Session.CommandRegistry.handlers[c.cmd.name].help()
			return err
		} else if _, found := findCMDElementWithName(s.Session.CommandRegistry.handlers[c.cmd.name].actions, action); !found {
			err = fmt.Errorf("invalid action for command '%s': %s", c.cmd.name, action)
			s.Session.CommandRegistry.handlers[c.cmd.name].help()
			return err
		}
		c.ctxValues["action"] = action
		return next(s, c)
	})
}

// =========== HANDLERS =============

func handlerExit(s *State, c *handlerContext) error {
	fmt.Println("Closing Pincher-CLI program...")
	*s.DoneChan <- true
	// NOTE: No actual error. This handler hijacks the error handling within
	// the repl to tell it that it should stop any looping.
	// The reason this is implemented is because without it, the loop will
	// still run and print out the REPL input prompt while waiting on
	// main() to close out the program.
	// TODO: find a more elegant way of doing this.
	// os.Exit() CAN'T be the solution, as the deferred Quit() function
	// under main.go would then not be called.
	return fmt.Errorf("HIJACK:EXIT")
}

func handlerClear(s *State, c *handlerContext) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// handlerHelp attempts to output the most relevant information possible to the user.
// If the user inquires simply for help, a list of all registered commands will be output.
// If asked for help with a specific command, its usage, actions, and commands will be output.
// If asked for help with one of a command's actions, its exclusive cmdElement usage will be output.
func handlerHelp(s *State, c *handlerContext) error {
	commandInquiry, _ := c.args.pfx()
	if handler, exists := s.Session.CommandRegistry.exists(commandInquiry); exists {
		actionInquiry, _ := c.args.pfx()
		if actionInquiry != "" {
			if cmdElement, found := findCMDElementWithName(handler.actions, actionInquiry); found {
				fmt.Println("USAGE: " + cmdElement.usage(true))
				return nil
			}
		}
		handler.help()
		return nil
	}
	if handler, exists := s.Session.CommandRegistry.exists("help"); exists {
		handler.help()
		fmt.Println("AVAILABLE COMMANDS: ")
		c.args.trackOptArgs(&c.cmd, "verbose")
		verbose, _ := c.args.pfx()
		registered := s.Session.CommandRegistry.GetRegisteredHandlers(verbose == "SET")
		sort.Slice(registered, func(i, j int) bool {
			return registered[i].priority < registered[j].priority
		})
		maxLen := MaxOfStrings(ExtractStrings(registered, func(c *cmdHandler) string { return c.name }))
		for _, handler := range registered {
			fmt.Printf("  %-*s  %s\n", maxLen, handler.name, handler.description)
		}

		return nil
	}
	return fmt.Errorf("could not get help for command: 'help'")
}

func handlerReady(s *State, c *handlerContext) error {
	isReady, err := s.Client.GetServerReady()
	if err != nil {
		return fmt.Errorf("server could not be reached; %w", err)
	}
	if isReady {
		fmt.Println("Server is ready!")
		return nil
	}
	fmt.Println("Server not ready.")
	return nil
}
