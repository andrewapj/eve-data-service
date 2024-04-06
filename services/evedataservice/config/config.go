package config

import (
	"github.com/andrewapj/dotenvconfig"
	"io/fs"
	"os"
	"sync"
)

const (
	defaultConfigPath = "local.env"
)

var mutex sync.Mutex

// Load loads the config for the application.
func Load(fSys fs.FS) {

	mutex.Lock()
	defer mutex.Unlock()

	path, ok := os.LookupEnv(ConfigPathKey())
	if !ok {
		path = defaultConfigPath
	}

	err := dotenvconfig.Load(fSys, path, dotenvconfig.Options{
		JsonLogging:    true,
		LoggingEnabled: true,
	})

	if err != nil {
		panic(err.Error())
	}
}
