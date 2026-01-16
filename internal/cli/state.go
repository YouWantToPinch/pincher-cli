package cli

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
)

// State represents the full state of the CLI during a session,
// regardless of user login. It knows the full context of all
// systems connected.
type State struct {
	DoneChan *chan bool
	Logger   *Logger
	Config   *config.Config
	Client   *client.Client
	Session  *cliSession
	CmdQueue chan<- string
	styles   *styles
}

func (s *State) NewSession() {
	s.Session = &cliSession{}
	s.Session.Init()
}

// GetPrompt returns the proper input
// prompt to print to the command line,
// dependent on user activity during
// this session.
func (s *State) GetPrompt() string {
	if s.styles != nil {
		return s.getStyledPrompt()
	}
	return s.getUnstyledPrompt()
}

func (s *State) getDiv(styled bool) string {
	div := "___________________________"
	if styled && s.styles != nil {
		return s.styles.Green.Render(div)
	}
	return div
}

// getUnstyledPrompt returns the pincher-cli prompt without pretty colors :(
func (s *State) getUnstyledPrompt() string {
	if s.Session.Username == "" {
		return "pin¢her > "
	} else if s.Client.ViewedBudget.Name == "" {
		return fmt.Sprintf("p¢/%s > ", s.Session.Username)
	} else {
		return fmt.Sprintf("p¢/%s[%s] > ", s.Session.Username, s.Client.ViewedBudget.Name)
	}
}

// getStyledPrompt returns the pincher-cli prompt with pretty colors :)
func (s *State) getStyledPrompt() string {
	cent := s.styles.Orange.Render("¢")
	chev := s.styles.White.Render(" > ")

	if s.Session.Username == "" {
		pin := s.styles.White.Render("pin")
		her := s.styles.White.Render("her")
		return pin + cent + her + chev
	} else {
		p := s.styles.White.Render("p")
		slash := s.styles.White.Render("/")

		if s.Client.ViewedBudget.Name == "" {
			return p + cent + slash + s.styles.White.Render(s.Session.Username) + chev
		} else {
			lbr := s.styles.Green.Render("[")
			rbr := s.styles.Green.Render("]")
			return p + cent + slash + s.styles.White.Render(s.Session.Username) + lbr + s.styles.White.Render(s.Client.ViewedBudget.Name) + rbr + chev
		}
	}
}

// cliSession represents the state of the CLI in regard to
// a logged-in user.
type cliSession struct {
	CommandRegistry *commandRegistry
	client.User
}

// Init preregisters all commands to the internal command registry.
func (s *cliSession) Init() {
	// initialize maps under command registry
	s.CommandRegistry = &commandRegistry{
		handlers: make(map[string]*cmdHandler),
		registry: make(map[string]registrationStatus),
	}
	// preregister ALL commands
	s.CommandRegistry.batchRegistration(makeBaseCommandHandlers(), Preregistered)
	s.CommandRegistry.preregister(makeBudgetCommandHandler())
	s.CommandRegistry.batchRegistration(makeResourceCommandHandlers(), Preregistered)
	// register base commands
	s.CommandRegistry.batchRegistration(makeBaseCommandHandlers(), Registered)
}

func (s *cliSession) OnLogin(user client.User) {
	// register commands that require login
	s.CommandRegistry.register("budget")
	s.User = user
	fmt.Printf("Logged in as user: %s\n", s.Username)
}

func (s *cliSession) OnViewBudget() {
	// register commands that require viewing a budget
	s.CommandRegistry.batchRegistration(makeResourceCommandHandlers(), Registered)
}

func (s *cliSession) OnLogout() {
	// deregister commands
	s.User = client.User{}
	s.CommandRegistry.deregisterNonBaseCommands()
}
