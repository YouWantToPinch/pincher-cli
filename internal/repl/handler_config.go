package repl

import (
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/config"
	"github.com/YouWantToPinch/pincher-cli/internal/tmodels"
	tea "github.com/charmbracelet/bubbletea"
)

func handlerConfig(s *State, cmd command) error {
	if _, res := cmd.hasOpt("edit"); res {
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
				s.Config = &newConfig
				s.Config.WriteToFile()
				fmt.Println("Saved configuration changes.")
			}
		}
		return nil
	} else if _, res := cmd.hasOpt("load"); res { // change _ to 'opt' for multi-config functionality
		userConfig := config.Config{}
		var err error
		userConfig, err = config.Read()
		if err != nil {
			s.Config = &userConfig
			return fmt.Errorf("Trouble loading config: %s", err.Error())
		}
		s.Config = &userConfig
		// TODO: Add multi-config functionality;
		// then print NAME of config profile here with opt[0].
		// fmt.Printf("Loaded configuration settings: %s\n", opt[0])
		fmt.Println("Loaded configuration settings.")
		return nil
	} else {
		return fmt.Errorf("Expected one of two options: ( --edit | --load)")
	}
}
