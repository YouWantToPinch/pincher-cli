package repl

import (
	"fmt"
	"os"
	"sort"
)

// ========= MIDDLEWARE ==============
// NOTE: This is handled by commandRegistry.run() instead!
// Keep this comment in case middleware written in this fashion
// may ever become necessary in the future.
/*
func middlewareAddArgsContext(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(s *State, c handlerContext) error {
		c.args.init(&c.cmd)
		return next(s, c)
	})
}
*/
// =========== HANDLERS =============

func handlerExit(s *State, c handlerContext) error {
	fmt.Println("Closing Pincher-CLI program...")
	os.Exit(0)
	return nil
}

func handlerHelp(s *State, c handlerContext) error {
	commandInquiry := c.args.pfx()
	if handler, exists := s.CommandRegistry.exists(commandInquiry); exists {
		handler.help()
		return nil
	}
	if handler, exists := s.CommandRegistry.exists("help"); exists {
		handler.help()
		fmt.Println("AVAILABLE COMMANDS: ")
		registered := s.CommandRegistry.GetRegisteredHandlers()
		sort.Slice(registered, func(i, j int) bool {
			return registered[i].priority < registered[j].priority
		})
		maxLen := MaxOfStrings(ExtractStrings(registered, func(c cmdHandler) string { return c.name }))
		for _, handler := range registered {
			fmt.Printf("  %-*s  %s\n", maxLen, handler.name, handler.description)
		}

		return nil
	}
	return fmt.Errorf("ERROR: Could not get help for command: 'help'")
}

func handlerLog(s *State, c handlerContext) error {
	return fmt.Errorf("ERROR: Command not implemented")
}

func handlerConnect(s *State, c handlerContext) error {
	s.Client.BaseUrl = s.Config.BaseURL
	fmt.Println("Set URL from config: " + s.Config.BaseURL)
	return nil
}

func handlerReady(s *State, c handlerContext) error {
	isReady, err := s.Client.GetServerReady()
	if err != nil {
		return fmt.Errorf("ERROR: Server could not be reached; %s", err)
	}
	if isReady {
		fmt.Println("Server is ready!")
		return nil
	}
	fmt.Println("Server not ready.")
	return nil
}

func handlerList(s *State, c handlerContext) error {
	return fmt.Errorf("ERROR: Command not implemented")
}

func handlerReport(s *State, c handlerContext) error {
	return fmt.Errorf("ERROR: Command not implemented")
}
