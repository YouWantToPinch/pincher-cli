package cli

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/config"
	ui "github.com/bntrtm/gostructui"
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
			return fmt.Errorf("action not implemented")
		}
	} else {
		return fmt.Errorf("action was not saved to context")
	}
}

func handleConfigEdit(s *State, c *handlerContext) error {
	customMenuSettings := &ui.MenuSettings{}
	customMenuSettings.Init()
	customMenuSettings.Header = "Edit your local configuration: "

	newConfig := s.Config.ConfigSettings
	configEditMenu, err := ui.InitialTModelStructMenu(&newConfig, []string{"RefreshToken"}, true, customMenuSettings)
	if err != nil {
		return err
	}
	p := tea.NewProgram(configEditMenu)
	if entry, err := p.Run(); err != nil {
		return err
	} else {
		if entry.(ui.TModelStructMenu).QuitWithCancel {
			fmt.Printf("Canceled user configuration changes.\n")
		} else {
			err = entry.(ui.TModelStructMenu).ParseStruct(&newConfig)
			if err != nil {
				return err
			}
			s.Config.ConfigSettings = newConfig
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
	userConfig, err := config.ReadFromFile()
	if err != nil {
		return fmt.Errorf("trouble loading config: %w", err)
	}
	s.Config = userConfig
	if s.Client.BaseURL != s.Config.BaseURL {
		s.Client.BaseURL = s.Config.BaseURL
		fmt.Println("Set URL from config: " + s.Config.BaseURL)
	}
	fmt.Println("Loaded configuration settings.")
	return nil
}
