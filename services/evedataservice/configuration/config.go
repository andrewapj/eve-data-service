package configuration

import (
	config "github.com/andrewapj/dotenvconfig"
	"io/fs"
	"os"
	"time"
)

const (
	defaultConfigPath = "local.env"
)

// Load loads the configuration for the application.
func Load(fSys fs.FS) {

	path, ok := os.LookupEnv(ConfigPathKey)
	if !ok {
		path = defaultConfigPath
	}

	err := config.Load(fSys, path, config.Options{
		JsonLogging:    true,
		LoggingEnabled: true,
	})

	if err != nil {
		panic(err.Error())
	}

	// Ensure the application runs in UTC
	time.Local = time.UTC
}
