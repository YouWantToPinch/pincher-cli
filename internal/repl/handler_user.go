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
		case "update":
			return handleUserUpdate(s, c)
		case "delete":
			// return handleUserDelete(s, c)
			fallthrough
		default:
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handleUserAdd(s *State, c *handlerContext) error {
	username, _ := c.args.pfx()
	password, _ := c.args.pfx()
	retypedPassword, _ := c.args.pfx()

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
	username, _ := c.args.pfx()
	password, _ := c.args.pfx()

	user, err := s.Client.LoginUser(username, password)
	if err != nil {
		return err
	}

	s.Client.LoggedInUser.Token = user.Token
	s.Client.LoggedInUser.RefreshToken = user.RefreshToken
	s.Client.LoggedInUser.User = user.User
	registerBudgetCommand(s, false)
	fmt.Printf("LOGGED IN as user: %s\n", user.Username)
	return nil
}

func handleUserUpdate(s *State, c *handlerContext) error {
	username, _ := c.args.pfx()
	password, _ := c.args.pfx()
	c.args.trackOptArgs(&c.cmd, "username")
	newUsername, err := c.args.pfx()
	if err != nil {
		newUsername = username
	}
	c.args.trackOptArgs(&c.cmd, "password")
	newPassword, err := c.args.pfx()
	if err != nil {
		newPassword = password
	} else {
		retypedNewPassword, _ := c.args.pfx()
		if newPassword != retypedNewPassword {
			return fmt.Errorf("fields for new password did not match")
		}
	}

	err = s.Client.UpdateUser(newUsername, newPassword)
	if err != nil {
		return err
	}

	fmt.Println("User updated with new information")
	return nil
}
