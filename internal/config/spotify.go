package config

type SpotifyConfig struct {
	ClientId     string `toml:"client-id"`
	ClientSecret string `toml:"client-secret"`
}
