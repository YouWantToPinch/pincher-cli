package repl

import (
	"fmt"
	"os"

	"github.com/YouWantToPinch/pincher-cli/internal/config"
	"github.com/YouWantToPinch/pincher-cli/internal/tmodels"
	tea "github.com/charmbracelet/bubbletea"
)

func handlerExit(s *State, cmd command) error {
	fmt.Println("Closing Pincher-CLI program...")
	os.Exit(0)
	return nil
}

func handlerHelp(s *State, cmd command) error {
	fmt.Println("Exit, Help, Config, Log, Connect, Report, Add, List")
	return nil
}

func handlerConfig(s *State, cmd command) error {
	switch cmd.args[0] {
	case "edit":
		fmt.Println("Edit your local configuration: ")
		newConfig := config.Config{}
		newConfig = *s.Config
		tmodel, err := tmodels.InitialModelMakeStruct(&newConfig, nil, true)
		if err != nil {
			return err
		}
		p := tea.NewProgram(tmodel)
		if entry, err := p.Run(); err != nil {
			return err
		} else {
			if entry.(tmodels.ModelMakeStruct).QuitWithCancel {
				fmt.Println(fmt.Sprintf("Canceled user configuration changes."))
			} else {
				err = entry.(tmodels.ModelMakeStruct).ParseStruct(&newConfig)
				s.Config = &newConfig
				s.Config.WriteToFile()
				fmt.Println(fmt.Sprintf("Saved configuration changes."))
			}
		}
		return nil
	case "load":
		userConfig := config.Config{}
		var err error
		userConfig, err = config.Read()
		if err != nil {
			s.Config = &userConfig
			return fmt.Errorf("Trouble loading config: %s", err.Error())
		}
		s.Config = &userConfig
		fmt.Printf("Loaded configuration settings: %s\n", cmd.args[0])
		return nil
	default:
		return fmt.Errorf("No valid arguments")
	}
}

func handlerLog(s *State, cmd command) error {
	return fmt.Errorf("ERROR: Command not implemented.")
}

func handlerConnect(s *State, cmd command) error {
	s.Client.BaseUrl = s.Config.BaseURL
	fmt.Println("Set URL from config: " + s.Config.BaseURL)
	return nil
}

func handlerReady(s *State, cmd command) error {
	isReady, err := s.Client.GetServerReady()
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	if isReady {
		fmt.Println("Server is ready!")
		return nil
	}
	fmt.Println("Server not ready.")
	return nil
}

func handlerList(s *State, cmd command) error {
	return fmt.Errorf("ERROR: Command not implemented.")
}

func handlerAdd(s *State, cmd command) error {
	if err := cmd.require(1); err != nil {
		return err
	}

	switch cmd.args[0] {
	case "user":
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
			//fmt.Println("For help with logging in, see: `help login`")
			return nil
		} else {
			return fmt.Errorf("ERROR: username already exists.")
		}
	default:
		return fmt.Errorf("ERROR: Command not implemented.")
	}
}

func handlerReport(s *State, cmd command) error {
	return fmt.Errorf("ERROR: Command not implemented.")
}

func handlerLogin(s *State, cmd command) error {
	if err := cmd.require(3); err != nil {
		return err
	}
	if cmd.args[0] != "user" {
		return fmt.Errorf("Usage: login user")
	}
	user, err := s.Client.LoginUser(cmd.args[1], cmd.args[2])
	if err != nil {
		return fmt.Errorf("ERROR: %s", err)
	}
	s.Client.JSONWebToken = user.Token
	fmt.Printf("Logged in as %s, using new access token: %s\n", user.Username, user.Token)
	return nil
}
