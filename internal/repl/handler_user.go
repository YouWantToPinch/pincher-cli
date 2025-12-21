package repl

import (
	"fmt"
)

func handlerUser(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "add":
			return handleUserAdd(s, c)
		case "login":
			return handleUserLogin(s, c)
		default:
			return fmt.Errorf("ERROR: action not implemented")
		}
	} else {
		return fmt.Errorf("ERROR: action was not saved to context")
	}
}

func handleUserAdd(s *State, c *handlerContext) error {
	username, err := c.args.pfx()
	if err != nil {
		return fmt.Errorf("missing argument for command: <username>")
	}
	password, err := c.args.pfx()
	if err != nil {
		return fmt.Errorf("missing argument for command: <password>")
	}
	retypedPassword, err := c.args.pfx()
	if err != nil {
		return fmt.Errorf("missing argument for command: <retyped password>")
	}

	if password != retypedPassword {
		return fmt.Errorf("ERROR: password fields did not match")
	}
	userCreated, err := s.Client.CreateUser(username, password)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	if userCreated {
		fmt.Println("User " + username + " successfully created with new password.")
		fmt.Println("For help logging in, see: `help user login`")
		return nil
	} else {
		return fmt.Errorf("ERROR: username already exists")
	}
}

func handleUserLogin(s *State, c *handlerContext) error {
	username, err := c.args.pfx()
	if err != nil {
		return fmt.Errorf("missing argument for command: <username>")
	}
	password, err := c.args.pfx()
	if err != nil {
		return fmt.Errorf("missing argument for command: <password>")
	}

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
	s.CommandRegistry.register("budget", &cmdHandler{
		cmdElement: cmdElement{
			name:        "budget",
			description: "Manage " + s.Client.LoggedInUser.Username + "'s budgets",
			priority:    100,
		},
		callback: middlewareValidateAction(handlerBudget),
		actions: []cmdElement{
			{
				name:      "add",
				arguments: []string{"name"},
				options: []cmdElement{
					{
						name:        "notes",
						description: "Give your budget some notes",
						arguments:   []string{"notes_value"},
					},
				},
			},
			{
				name: "list",
				options: []cmdElement{
					{
						name:        "role",
						description: "Filter results by user role. Can be ADMIN, MANAGER, CONTRIBUTOR, or VIEWER.",
						arguments:   []string{"role_title"},
					},
				},
			},
			{
				name:      "view",
				arguments: []string{"budget_name"},
			},
		},
	})

	return nil
}
