package player

import (
	"context"
	"log"
	"strings"

	"github.com/zmb3/spotify/v2"
)

type SpotifyPlayer struct {
	client    *spotify.Client
	isPlaying bool
}

// Impl Player interface
func (s *SpotifyPlayer) IsPlaying() bool {
	return s.isPlaying
}

func (s *SpotifyPlayer) CurrentSong() Song {
	playerState := s.playerState()
	if playerState.Item == nil {
		log.Printf("Unable to get player state: %v", playerState)
		return Song{}
	}

	var artistNames []string
	for _, artist := range playerState.Item.Artists {
		artistNames = append(artistNames, artist.Name)
	}

	// Tracks from a podcast come through with Progress == Duration. This is a hack around that.
	// Obviously we're beat if we're in the middle of a track
	progress := playerState.Progress
	if progress == playerState.Item.Duration {
		progress = 0
	}

	return Song{
		Name:     playerState.Item.Name,
		Artist:   strings.Join(artistNames, ", "),
		Length:   playerState.Item.Duration,
		Position: progress,
	}
}

func (s *SpotifyPlayer) PlayPause() (Song, error) {
	ctx := context.Background()
	ps := s.playerState()
	var err error
	if ps.Playing {
		err = s.client.Pause(ctx)
	} else {
		err = s.client.Play(ctx)
	}

	if err == nil {
		s.isPlaying = !s.isPlaying
	}

	return s.CurrentSong(), err
}

func (s *SpotifyPlayer) NextSong() (Song, error) {
	ctx := context.Background()
	err := s.client.Next(ctx)

	return s.CurrentSong(), err
}

func (s *SpotifyPlayer) PreviousSong() (Song, error) {
	ctx := context.Background()

	err := s.client.Previous(ctx)
	return s.CurrentSong(), err
}

func (s *SpotifyPlayer) playerState() *spotify.PlayerState {
	playerState, err := s.client.PlayerState(context.Background())
	if err != nil {
		log.Fatalf("Unable to get player state: %s", err)
	}

	// always sync play/pause on query
	s.isPlaying = playerState.Playing
	return playerState
}

func New(client *spotify.Client) *SpotifyPlayer {
	return &SpotifyPlayer{client, false}
}
