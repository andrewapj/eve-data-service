package testhelper

import (
	"fmt"
	"github.com/andrewapj/arcturus/config"
	"os"
	"path/filepath"
)

const testConfigFile = "local-test.env"

// SetTestConfig sets the config for tests.
func SetTestConfig() {

	dir, err := GetRootDir()
	if err != nil {
		panic(err.Error())
	}

	err = os.Setenv(config.ConfigPathKey(), testConfigFile)
	if err != nil {
		panic(err.Error())
	}

	config.Load(os.DirFS(dir))
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
