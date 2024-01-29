// Handler for loading various providers (Spotify, Soundcloud, etc.)
package providers

import (
	"log"
	"tmix/internal/config"
	cfg "tmix/internal/config"
	"tmix/internal/player"
)

type Song struct {
	ProviderId string
	Name       string
	Artist     string
}

type Playlist struct {
	ProviderId  string
	Name        string
	Description string
}

type AbstractMusicProvider struct {
	loggedIn bool
}

func (m AbstractMusicProvider) LoggedIn() bool {
	return m.loggedIn
}

type MusicProvider interface {
	Player() player.Player
	Name() string
	LoggedIn() bool
	Login()
	FetchPlaylists() []Playlist
}

func LoadProviders(config *config.ProviderConfig) []MusicProvider {
	spot := NewSpotify(config.Spotify)
	log.Printf("Loading new spotify provider with: %v", config.Spotify)
	spot.cache = cfg.New(config.AuthTokenCacheDir)
	return []MusicProvider{spot}
}
