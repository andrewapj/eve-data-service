package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// SetTestConfig sets the config for tests.
func SetTestConfig() {

	dir, err := GetRootDir()
	if err != nil {
		panic(err.Error())
	}

	err = os.Setenv(ConfigPathKey(), "local-test.env")
	if err != nil {
		panic(err.Error())
	}

	Load(os.DirFS(dir))
}

// GetRootDir finds the root application path where 'main.go' is located.
// This allows tests run in an IDE to find the config files.
func GetRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		mainGoPath := filepath.Join(dir, "main.go")
		if _, err := os.Stat(mainGoPath); err == nil {
			return dir, nil
		} else if os.IsNotExist(err) {
			parentDir := filepath.Dir(dir)
			if parentDir == dir {
				break
			}
			dir = parentDir
		} else {
			return "", fmt.Errorf("error checking for main.go: %v", err)
		}
	}

	return "", fmt.Errorf("main.go not found")
}
