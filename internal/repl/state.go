package repl

import (
	"encoding/json"
	"fmt"

	"github.com/YouWantToPinch/pincher-cli/internal/cache"
	"github.com/YouWantToPinch/pincher-cli/internal/client"
	"github.com/YouWantToPinch/pincher-cli/internal/config"
	file "github.com/YouWantToPinch/pincher-cli/internal/filemgr"
)

type State struct {
	DoneChan        *chan bool
	Logger          *Logger
	Config          *config.Config
	Client          *client.Client
	CommandRegistry *commandRegistry
}

// LoadCache loads a previous session into memory,
// if it exists.
func (s *State) LoadCache() error {
	cachePath, err := file.GetCacheFilepath("cache.json")
	if err != nil {
		return fmt.Errorf("failed to load cache: %w", err)
	}

	loadedCache, err := file.ReadJSONFromFile[cache.Cache](cachePath)
	if err != nil {
		return err
	}

	s.Client.Cache.Set(loadedCache.CachedEntries)

	if s.Config.StayLoggedIn {
		data, found := s.Client.Cache.Get("logged_in_user")
		if !found {
			return fmt.Errorf("failed to load cache: %w", err)
		}
		err := json.Unmarshal(data, &s.Client.LoggedInUser)
		if err != nil {
			return fmt.Errorf("failed to load cache: %w", err)
		}
	}

	return nil
}

// SaveCache saves the current session to a local *.json file
// under the pincher-cli cache directory.
func (s *State) SaveCache() error {
	if s.Config.StayLoggedIn {
		s.Client.Cache.Delete("logged_in_user")
		userBytes, err := json.Marshal(s.Client.LoggedInUser)
		if err != nil {
			return fmt.Errorf("failed to save cache: %w", err)
		}
		s.Client.Cache.Add("logged_in_user", userBytes, true)
	}

	cachePath, err := file.GetCacheFilepath("cache.json")
	if err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	err = file.WriteAsJSON(s.Client.Cache, cachePath)
	if err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}
