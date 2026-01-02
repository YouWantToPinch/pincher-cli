// Package config handles interpretation of pincher-cli user configuration
package config

import (
	"encoding/json"
	"os"

	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
)

type Config struct {
	BaseURL        string `json:"db_url"`
	VimKeysEnabled bool   `json:"vim_keys_enabled"`
}

func (c *Config) New(dbURL, username string) error {
	newCfg := Config{
		BaseURL:        dbURL,
		VimKeysEnabled: true,
	}
	err := newCfg.WriteToFile()
	if err != nil {
		return err
	}

	return nil
}

func ReadFromFile() (*Config, error) {
	confPath, err := file.GetConfigFilepath("cli.conf")
	if err != nil {
		return nil, err
	}

	config, err := file.ReadJSONFromFile[Config](confPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) WriteToFile() error {
	path, err := file.GetConfigFilepath("cli.conf")
	if err != nil {
		return err
	}

	err = file.WriteAsJSON(c, path)
	if err != nil {
		return err
	}

	return nil
}
