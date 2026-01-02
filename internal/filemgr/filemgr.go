// Package filemgr is used to get and interact with
// directories and files pertinent to running the CLI.
package filemgr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type userDirFunc func() (string, error)

// GetFilePath gets the path of a specified file under the specified base user directory.
// Pincher aims to adhere to XDG Base Directory Specification.
// For the given each base directory, GetConfigPath will look within
// the ordered subdirectories as specified for the file.
func getFilePath(fn userDirFunc, subdirs []string, filename string) (string, error) {
	userBaseDir, err := fn()
	if err != nil {
		return "", err
	}
	path := filepath.Join(append([]string{userBaseDir}, subdirs...)...)
	return filepath.Join(path, filename), nil
}

// GetConfigFilepath returns the path of a specific config file under the application's config directory.
func GetConfigFilepath(filename string) (string, error) {
	return getFilePath(os.UserConfigDir, []string{"pincher"}, filename)
}

// GetLogFilepath returns the path of a specific log file under the application's "logs" directory.
func GetLogFilepath(filename string) (string, error) {
	return getFilePath(os.UserHomeDir, []string{".local", "share", "pincher", "logs"}, filename)
}

// GetCacheFilepath returns the path of a specific cache file under the application's cache directory.
func GetCacheFilepath(filename string) (string, error) {
	return getFilePath(os.UserCacheDir, []string{"pincher"}, filename)
}

// WriteAsJSON writes a given struct with JSON tags
// as JSON to a specified filepath
func WriteAsJSON(dataStruct any, filePath string) error {
	jsonData, err := json.MarshalIndent(dataStruct, "", " \t")
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0o666)
	if err != nil {
		return err
	}
	return nil
}

func ReadJSONFromFile[T any](filepath string) (*T, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var t T
	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
