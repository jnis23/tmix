package config

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

var (
	CacheDir  = loadCacheDir()
	ConfigDir = loadConfigDir()
	LogPath   = fmt.Sprintf("%s/debug.log", loadCacheDir())
)

type ProviderConfig struct {
	Spotify           SpotifyConfig `toml:"spotify"`
	AuthTokenCacheDir string        `toml:"auth-token-cache-dir"`
}

type Config struct {
	Providers *ProviderConfig `toml:"providers"`
	// This is only applicable for async players like Spotify where the "Player" is behind an API
	SyncRate int `toml:"sync-rate,commented" comment:"Rate (in seconds) at which we sync the playbar with the player. If 0, we disable sync and will manually increment every second until the song is over."`
}

func homeOrCurrentDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("Can't access user Home dir. Using current dir for cache...")
		wd, err := os.Getwd()
		if err != nil {
			fmt.Print("Unable to get home dir or current working directory. Shutting down...")
			os.Exit(1)
		}
		return wd
	}
	return home
}

func loadCacheDir() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		fmt.Print("Unable to access user cache dir. Trying User home dir...")
		cache = homeOrCurrentDir()
	}

	path := fmt.Sprintf("%s/tmix", cache)
	checkCreateDir(path)
	return path
}

func loadConfigDir() string {
	config := fmt.Sprintf("%s/.config/tmix", homeOrCurrentDir())

	checkCreateDir(config)
	return config
}

func checkCreateDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			log.Print("Creating new dir at ", dir)
			err = os.Mkdir(dir, fs.ModePerm)
			if err != nil {
				log.Fatalf("Unable to create dir at %s. %s", dir, err)
			}
		}
	}
}

func LoadConfig(configPath *string) *Config {
	var path string
	if configPath == nil || *configPath == "" {
		path = fmt.Sprintf("%s/config.toml", ConfigDir)
	} else {
		path = *configPath
	}
	// Try to load config from expected location
	cfg, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	parsed := &Config{}

	err = toml.Unmarshal(cfg, parsed)

	if err != nil {
		log.Fatal("Error while loading config: ", err)
	}

	return parsed
}
