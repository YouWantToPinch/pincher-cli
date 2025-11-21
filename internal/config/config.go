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

func (c *Config) New(dbUrl, username string) error {
	newCfg := Config{
		BaseURL:        dbUrl,
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
		return fmt.Errorf("Error: %v", err)
	}
	err = os.WriteFile(getConfigPath(), jsonData, 0666)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	return nil
}
