// Package filemgr is used to get and interact with
// directories and files pertinent to running the CLI.
package filemgr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// get the path of a specified file under the application "logs" directory
func GetLogPath(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".local", "share", "pincher", "logs", filename), nil
}

// get the path of a specified file under the application ".config/pincher" directory
func GetConfigPath(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "pincher", filename), nil
}

// WriteAsJSON writes a given struct with JSON tags
// as JSON to a specified filepath
func WriteAsJSON(dataStruct any, filepath string) error {
	jsonData, err := json.MarshalIndent(dataStruct, "", " \t")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, jsonData, 0o666)
	if err != nil {
		return err
	}
	return nil
}
