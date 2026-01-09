package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/YouWantToPinch/pincher-cli/internal/cli"
	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
)

// Quit is used by main() to ensure deletion of empty log files
func Quit(logger *cli.Logger) {
	err := logger.Close()
	if err != nil {
		fmt.Printf("LOGGER ERROR: %s\n", err)
	}
}

func main() {
	done := make(chan bool)

	cliState := &cli.State{DoneChan: &done}
	cfg, err := config.ReadFromFile()
	if err != nil {
		fmt.Printf("CONFIG ERROR: %s (does a config file exist?)\n", err)
	}
	cliState.Config = cfg

	cliState.Logger = &cli.Logger{}
	err = cliState.Logger.New(slog.LevelInfo)
	if err != nil {
		fmt.Printf("LOGGER ERROR: %s\n", err)
	}
	defer Quit(cliState.Logger)

	client := client.NewClient(time.Second*10, time.Minute*5, cliState.Config.BaseURL)
	cliState.Client = &client

	// give the client the stored refresh token so it will load cache
	if cliState.Config.StayLoggedIn {
		cliState.Client.RefreshToken = cliState.Config.RefreshToken
		if err := cliState.Client.LoadCacheFile(); err != nil {
			fmt.Printf("CACHE ERROR: %s\n", err)
		}
	} else {
		cliState.Config.RefreshToken = ""
	}

	// run the repl until it is closed from within
	go func() {
		cli.StartRepl(cliState)
	}()

	<-done
	if cliState.Config.StayLoggedIn {
		cliState.Config.RefreshToken = cliState.Client.RefreshToken
		err := cliState.Config.WriteToFile()
		if err != nil {
			fmt.Printf("CONFIG ERROR: %s\n", err)
		}
		err = cliState.Client.SaveCacheFile()
		if err != nil {
			fmt.Printf("CACHE ERROR: %s\n", err)
		} else {
			cliState.Config.RefreshToken = ""
		}
	}
}
