package config

import (
	"fmt"
	"log"
	"os"
	"tmix/internal/tui"

	toml "github.com/pelletier/go-toml/v2"
)

func LoadConfig(configPath *string) *tui.Config {
	var path string
	if configPath == nil || *configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("No config file set and unable to access user home directory: ", err)
		}

		path = fmt.Sprintf("%s/.config/tmix/config.toml", home)

	} else {
		path = *configPath
	}
	// Try to load config from expected location
	cfg, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	parsed := &tui.Config{}

	err = toml.Unmarshal(cfg, parsed)

	if err != nil {
		log.Fatal("Error while loading config: ", err)
	}

	return parsed
}
