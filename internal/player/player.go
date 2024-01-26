package player

type Song struct {
	Name   string
	Artist string
	// Length of the current song in milliseconds
	Length int
	// Current position in the song in milliseconds
	Position int
}

func (s *Song) Equals(o Song) bool {
	return s.Artist == o.Artist && s.Name == o.Name
}

type Player interface {
	IsPlaying() bool
	CurrentSong() Song
	PlayPause() (Song, error)
	NextSong() (Song, error)
	PreviousSong() (Song, error)
}
