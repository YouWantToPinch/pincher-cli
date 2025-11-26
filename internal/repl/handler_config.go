package repl

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/config"
	"github.com/YouWantToPinch/pincher-cli/internal/tmodels"
	tea "github.com/charmbracelet/bubbletea"
)

func handlerConfig(s *State, c handlerContext) error {
	action := c.args.pfx()
	switch action {
	case "edit":
		return handle_editConfig(s, c)

	case "load":
		return handle_LoadConfig(s, c)
	default:
		return fmt.Errorf("Expected one of two options: ( --edit | --load)")

	}
}

func handle_editConfig(s *State, c handlerContext) error {
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
			fmt.Printf("Canceled user configuration changes.\n")
		} else {
			err = entry.(tmodels.ModelMakeStruct).ParseStruct(&newConfig)
			if err != nil {
				return err
			}
			s.Config = &newConfig
			s.Config.WriteToFile()
			fmt.Println("Saved configuration changes.")
		}
		return nil
	}
}

func handle_LoadConfig(s *State, c handlerContext) error {
	userConfig := config.Config{}
	var err error
	userConfig, err = config.Read()
	if err != nil {
		s.Config = &userConfig
		return fmt.Errorf("Trouble loading config: %s", err.Error())
	}
	s.Config = &userConfig
	fmt.Println("Loaded configuration settings.")
	return nil
}
