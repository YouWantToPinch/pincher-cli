// Package config handles CLI settings and optionally stores refresh tokens
package config

import (
	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
)

type ConfigSettings struct {
	BaseURL         string `json:"db_url" smname:"Database URL" smdes:"URL of the server to connect to"`
	StayLoggedIn    bool   `json:"stay_logged_in" smname:"Stay Logged In" smdes:"Keep a login session alive on exit."`
	CurrencyISOCode string `json:"currency_iso_code" smname:"Currency ISO" smdes:"The ISO Code of the currency desired for monetary visualization"`
	VimKeysEnabled  bool   `json:"vim_keys_enabled" smname:"Vim Keys Enabled" smdes:"Use vim keys to navigate CLI menus."`
}

// Config represents a configuration specific to the local machine.
type Config struct {
	ConfigSettings
	RefreshToken string `json:"refresh_token"`
}

func (c *Config) NewConfigFile(dbURL string) error {
	c.SetDefaults(dbURL)
	err := c.WriteToFile()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SetDefaults(dbURL string) {
	c.ConfigSettings = ConfigSettings{
		BaseURL:         dbURL,
		StayLoggedIn:    true,
		CurrencyISOCode: "USD",
		VimKeysEnabled:  true,
	}
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
