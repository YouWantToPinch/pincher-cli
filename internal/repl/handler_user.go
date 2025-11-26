package repl

import (
	"fmt"
)

func handlerUser(s *State, c handlerContext) error {
	action := c.args.pfx()
	switch action {
	case "add":
		return handle_userAdd(s, c)
	case "login":
		return handle_userLogin(s, c)
	case "":
		return fmt.Errorf("ERROR: no action specified")
	default:
		return fmt.Errorf("ERROR: invalid action for user: %s", action)
	}
}

func handle_userAdd(s *State, c handlerContext) error {
	username := c.args.pfx()
	password := c.args.pfx()
	retypedPassword := c.args.pfx()

	if password != retypedPassword {
		return fmt.Errorf("ERROR: password fields did not match.")
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
		return fmt.Errorf("ERROR: username already exists.")
	}
}

func handle_userLogin(s *State, c handlerContext) error {
	username := c.args.pfx()
	password := c.args.pfx()

	user, err := s.Client.LoginUser(username, password)
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	s.Client.JSONWebToken = user.Token
	fmt.Printf("Logged in as %s, using new access token: %s\n", user.Username, user.Token)
	return nil
}
