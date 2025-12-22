package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
	"github.com/YouWantToPinch/pincher-cli/internal/repl"
)

// Quit is used by main() to ensure deletion of empty log files
func Quit(logger *repl.Logger) {
	err := logger.Close()
	if err != nil {
		fmt.Printf("LOGGER ERROR: %s\n", err.Error())
	}
}

func main() {
	done := make(chan bool)

	cliState := &repl.State{DoneChan: &done}
	cfg, err := config.ReadFromFile("cli.conf")
	if err != nil {
		fmt.Printf("CONFIG ERROR: %s\n", err.Error())
		fmt.Println("(Does a config file exist?)")
	}
	cliState.Config = &cfg

	cliState.Logger = &repl.Logger{}
	err = cliState.Logger.New(slog.LevelInfo)
	if err != nil {
		fmt.Printf("LOGGER ERROR: %s\n", err)
	}
	defer Quit(cliState.Logger)

	client := client.NewClient(time.Second*10, time.Minute*5, cliState.Config.BaseURL)
	cliState.Client = &client

	// run the repl until it is closed from within
	go func() {
		repl.StartRepl(cliState)
	}()

	<-done
}
