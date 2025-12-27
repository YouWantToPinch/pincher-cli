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
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
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
		return fmt.Errorf("password fields did not match")
	}
	userCreated, err := s.Client.CreateUser(username, password)
	if err != nil {
		return err
	}
	if userCreated {
		fmt.Println("User " + username + " successfully created with new password.")
		fmt.Println("For help logging in, see: `help user login`")
		return nil
	} else {
		return fmt.Errorf("username already exists")
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
		return err
	}
	s.Client.LoggedInUser.JSONWebToken = user.Token
	s.Client.LoggedInUser.Username = user.Username
	registerBudgetCommand(s, false)
	fmt.Printf("Logged in as %s, using new access token: %s\n", user.Username, user.Token)
	return nil
}
