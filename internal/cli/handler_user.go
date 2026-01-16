package cli

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
		case "logout":
			return handleUserLogout(s, c)
		case "update":
			return handleUserUpdate(s, c)
		case "delete":
			return handleUserDelete(s, c)
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
		fmt.Println("For help logging in, see: `help user -a login`")
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

	s.Session.OnLogin(user)

	c.args.trackOptArgs(&c.cmd, "view-budget")
	budgetToView, _ := c.args.pfx()
	if budgetToView != "" {
		s.CmdQueue <- "budget view " + fmt.Sprintf(`"%s"`, budgetToView)
	}

	return nil
}

func handleUserLogout(s *State, c *handlerContext) error {
	if s.Client.RefreshToken == "" {
		fmt.Println("No user logged in.")
		return nil
	}
	err := s.Client.RevokeRefreshToken()
	if err != nil {
		return err
	}
	s.Client.ClearCache()
	s.Session.OnLogout()

	fmt.Println("User logged out.")
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

func handleUserDelete(s *State, c *handlerContext) error {
	username, _ := c.args.pfx()
	password, _ := c.args.pfx()
	retypedPassword, _ := c.args.pfx()

	if password != retypedPassword {
		return fmt.Errorf("password fields did not match")
	}
	err := s.Client.DeleteUser(username, password)
	if err != nil {
		return err
	}

	fmt.Println("User " + username + " successfully deleted.")
	return nil
}
