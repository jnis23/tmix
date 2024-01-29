package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"

	"golang.org/x/oauth2"
)

type TokenCache struct {
	dir     string
	Enabled bool
}

func New(dir string) *TokenCache {
	var d = dir
	if dir == "" {
		d = CacheDir
	}

	cache := &TokenCache{d, true}
	cache.checkAndCreateFiles()
	return cache
}

func (t *TokenCache) tokenFileName() string {
	return fmt.Sprintf("%s/token", t.dir)
}

func (t *TokenCache) StoreToken(tok *oauth2.Token) error {
	f, err := os.Create(t.tokenFileName())
	defer f.Close()
	if err != nil {
		return err
	}
	j, err := json.Marshal(tok)
	if err != nil {
		return err
	}
	_, err = f.Write(j)

	return err
}

func (t *TokenCache) FetchToken() *oauth2.Token {
	f, err := os.ReadFile(t.tokenFileName())
	if err != nil {
		log.Fatalf("Unable to open token file: %s", err)
	}
	if len(f) == 0 {
		return nil
	}
	tok := oauth2.Token{}

	err = json.Unmarshal(f, &tok)
	if err != nil {
		log.Fatalf("Unable to get token from file: %s", err)
	}

	return &tok
}

func (t *TokenCache) checkAndCreateFiles() {
	if _, err := os.Stat(t.dir); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(t.dir, fs.ModePerm)
			log.Printf("Creating temp dir at %s", t.dir)
			if err != nil {
				log.Fatalf("Unable to create dir at %s", t.dir)
			}
		}
	}

	checkOrCreateFile(t.tokenFileName())
}

func checkOrCreateFile(filename string) {
	var err error
	if _, err = os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			_, err = os.Create(filename)
		}
	}
	if err != nil {
		log.Fatalf("Failed to create file %s: %s", filename, err)
	}
}
