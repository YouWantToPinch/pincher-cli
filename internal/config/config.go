// Package config handles interpretation of pincher-cli user configuration
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "pincherconfig.json"

type Config struct {
	BaseURL        string `json:"db_url"`
	VimKeysEnabled bool   `json:"vim_keys_enabled"`
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	return fmt.Sprintf("%s/.config/pincher/%s", homeDir, configFileName)
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

func Read() (Config, error) {
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		return Config{}, fmt.Errorf("ERROR: %s", err.Error())
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("ERROR: %s", err.Error())
	}
	return config, nil
}

func (c *Config) WriteToFile() error {
	jsonData, err := json.MarshalIndent(c, "", " \t")
	if err != nil {
		return err
	}
	err = os.WriteFile(getConfigPath(), jsonData, 0o666)
	if err != nil {
		return err
	}
	return nil
}
