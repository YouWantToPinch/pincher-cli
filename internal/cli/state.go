package cli

import (
	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
)

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

func (s *cliSession) OnLogin() {
	// register commands that require login
	s.CommandRegistry.register("budget")
}

func (s *cliSession) OnViewBudget() {
	// register commands that require viewing a budget
	s.CommandRegistry.batchRegistration(makeResourceCommandHandlers(), Registered)
}

func (s *cliSession) OnLogout() {
	// deregister commands
	s.CommandRegistry.deregisterNonBaseCommands()
}

func (s *cliSession) OnExit() {
	// deregister commands
	s.CommandRegistry.deregisterNonBaseCommands()
}

type State struct {
	DoneChan *chan bool
	Logger   *Logger
	Config   *config.Config
	Client   *client.Client
	Session  *cliSession
}

func (s *State) NewSession() {
	s.Session = &cliSession{}
	s.Session.Init()
}
