package repl

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/config"
	"github.com/YouWantToPinch/pincher-cli/internal/tmodels"
	tea "github.com/charmbracelet/bubbletea"
)

func handlerConfig(s *State, c *handlerContext) error {
	if val, ok := c.ctxValues["action"]; ok {
		switch val {
		case "edit":
			return handleConfigEdit(s, c)
		case "load":
			return handleConfigLoad(s, c)
		default:
			return fmt.Errorf("ERROR: action not implemented")
		}
	} else {
		return fmt.Errorf("ERROR: action was not saved to context")
	}
}

func handleConfigEdit(s *State, c *handlerContext) error {
	fmt.Println("Edit your local configuration: ")
	newConfig := *s.Config
	tmodel, err := tmodels.InitialTModelStructMenu(&newConfig, nil, true)
	if err != nil {
		return err
	}
	p := tea.NewProgram(tmodel)
	if entry, err := p.Run(); err != nil {
		return err
	} else {
		if entry.(tmodels.TModelStructMenu).QuitWithCancel {
			fmt.Printf("Canceled user configuration changes.\n")
		} else {
			err = entry.(tmodels.TModelStructMenu).ParseStruct(&newConfig)
			if err != nil {
				return err
			}
			s.Config = &newConfig
			err := s.Config.WriteToFile()
			if err != nil {
				return err
			}
			fmt.Println("Saved configuration changes.")
		}
		return nil
	}
}

func handleConfigLoad(s *State, c *handlerContext) error {
	var err error
	userConfig, err := config.Read()
	if err != nil {
		s.Config = &userConfig
		return fmt.Errorf("trouble loading config: %s", err.Error())
	}
	s.Config = &userConfig
	fmt.Println("Loaded configuration settings.")
	return nil
}
