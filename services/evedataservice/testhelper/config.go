package testhelper

import (
	"fmt"
	"github.com/andrewapj/arcturus/configuration"
	"os"
	"path/filepath"
	"time"
)

const testConfigFile = "local-test.env"

// SetTestConfig sets the configuration for tests.
func SetTestConfig() {
	dir, err := GetRootDir()
	if err != nil {
		panic(err.Error())
	}

	err = os.Setenv(configuration.ConfigPathKey, testConfigFile)
	if err != nil {
		panic(err.Error())
	}
	configuration.Load(os.DirFS(dir))

	// Ensure the application runs in UTC
	time.Local = time.UTC
}

// GetRootDir finds the root application path where 'main.go' is located.
// This allows tests run in an IDE to find the configuration files.
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
