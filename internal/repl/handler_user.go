package repl

import (
	"fmt"
)

func handlerUser(s *State, c handlerContext) error {
	action := c.args.pfx()
	switch action {
	case "add":
		return handleUserAdd(s, c)
	case "login":
		return handleUserLogin(s, c)
	case "":
		return fmt.Errorf("ERROR: no action specified")
	default:
		return fmt.Errorf("ERROR: invalid action for user: %s", action)
	}
}

func handleUserAdd(s *State, c handlerContext) error {
	username := c.args.pfx()
	password := c.args.pfx()
	retypedPassword := c.args.pfx()

	if password != retypedPassword {
		return fmt.Errorf("ERROR: password fields did not match")
	}
	userCreated, err := s.Client.CreateUser(username, password)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	if userCreated {
		fmt.Println("User " + username + " successfully created with new password.")
		fmt.Println("For help logging in, see: `help login`")
		return nil
	} else {
		return fmt.Errorf("ERROR: username already exists")
	}
}

func handleUserLogin(s *State, c handlerContext) error {
	username := c.args.pfx()
	password := c.args.pfx()

	user, err := s.Client.LoginUser(username, password)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	s.Client.LoggedInUser.JSONWebToken = user.Token
	s.Client.LoggedInUser.Username = user.Username
	err = registerResourceCommands(s)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	fmt.Printf("Logged in as %s, using new access token: %s\n", user.Username, user.Token)
	return nil
}

func registerResourceCommands(s *State) error {
	s.CommandRegistry.register("budget", cmdHandler{
		name:        "budget",
		description: "Manage " + s.Client.LoggedInUser.Username + "'s budgets",
		priority:    100,
		callback:    handlerBudget,
		opts: []cmdOption{
			{
				word:        "role",
				description: "Filter results by user role. Can be ADMIN, MANAGER, CONTRIBUTOR, or VIEWER.",
				argCount:    1,
			},
			{
				word:        "notes",
				description: "Give your budget some notes",
				argCount:    1,
			},
		},
		usage: `budget <action> (arguments) [options]
budget add <name> [ --notes ]
budget list [ --role ]
budget view <budget_name>`,
	})
	return nil
}
