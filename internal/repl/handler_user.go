package repl

import (
	"fmt"
)

func handlerUser(s *State, cmd command) error {
	if err := cmd.require(1); err != nil {
		return err
	}

	switch cmd.args[0] {
	case "add":
		if err := cmd.require(4); err != nil {
			return err
		}

		if cmd.args[2] != cmd.args[3] {
			return fmt.Errorf("ERROR: password fields did not match.")
		}
		userCreated, err := s.Client.CreateUser(cmd.args[1], cmd.args[2])
		if err != nil {
			return fmt.Errorf("ERROR: %s", err)
		}
		if userCreated {
			fmt.Println("User " + cmd.args[1] + " successfully created with new password.")
			fmt.Println("For help with logging in, see: `help login`")
			return nil
		} else {
			return fmt.Errorf("ERROR: username already exists.")
		}
	case "login":
		if err := cmd.require(3); err != nil {
			return err
		}

		user, err := s.Client.LoginUser(cmd.args[1], cmd.args[2])
		if err != nil {
			return fmt.Errorf("ERROR: %s", err)
		}
		s.Client.JSONWebToken = user.Token
		fmt.Printf("Logged in as %s, using new access token: %s\n", user.Username, user.Token)
		return nil
	default:
		return fmt.Errorf("ERROR: Command not implemented.")
	}
}

func handle_userAdd(s *State, cmd command) error {
	if err := cmd.require(4); err != nil {
		return err
	}

	if cmd.args[2] != cmd.args[3] {
		return fmt.Errorf("ERROR: password fields did not match.")
	}
	userCreated, err := s.Client.CreateUser(cmd.args[1], cmd.args[2])
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	if userCreated {
		fmt.Println("User " + cmd.args[1] + " successfully created with new password.")
		fmt.Println("For help with logging in, see: `help login`")
		return nil
	} else {
		return fmt.Errorf("ERROR: username already exists.")
	}
}
