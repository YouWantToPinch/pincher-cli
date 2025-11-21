package main

import (
	"fmt"
	"time"

	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
	"github.com/YouWantToPinch/pincher-cli/internal/repl"
)

func main() {
	cliState := repl.State{}
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		fmt.Println("(Does a config file exist?)")
	}
	cliState.Config = &cfg

	client := client.NewClient(time.Second*10, time.Minute*5, cliState.Config.BaseURL)
	cliState.Client = &client

	repl.StartRepl(&cliState)

}
