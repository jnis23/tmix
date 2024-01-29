package tui

import "tmix/internal/providers"

type Config struct {
	Providers *providers.ProviderConfig `toml:"providers"`
	// This is only applicable for async players like Spotify where the "Player" is behind an API
	SyncRate int `toml:"sync-rate,commented" comment:"Rate (in seconds) at which we sync the playbar with the player. If 0, we disable sync and will manually increment every second until the song is over."`
}
