package tui

import (
	"fmt"
	"log"
	"time"
	pl "tmix/internal/player"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type playerModel struct {
	song    pl.Song
	player  pl.Player
	playbar progress.Model
}

// Keys
type playerKeyMap struct {
	PlayPause key.Binding
	NextSong  key.Binding
	PrevSong  key.Binding
}

var playerKeys = playerKeyMap{
	PlayPause: key.NewBinding(
		key.WithKeys(" "),
	),
	NextSong: key.NewBinding(
		key.WithKeys("n"),
	),
	PrevSong: key.NewBinding(
		key.WithKeys("p"),
	),
}

type changedEvent int

const (
	None changedEvent = iota
	PlayPause
	PrevNext
)

// Commands
type songChangedMsg struct {
	song  pl.Song
	event changedEvent
}

type pbTickMsg time.Time

func (p *playerModel) SongChangedCmd() tea.Msg {
	return songChangedMsg{p.player.CurrentSong(), None}
}

func (p *playerModel) PlayPauseCmd() tea.Msg {
	s, err := p.player.PlayPause()
	if err != nil {
		log.Fatalf("Failed to play/pause: %s", err)
	}
	return songChangedMsg{s, PlayPause}
}

func (p *playerModel) NextSongCmd() tea.Msg {
	s, err := p.player.NextSong()
	if err != nil {
		log.Fatalf("Failed to get next song: %s", err)
	}
	return songChangedMsg{s, PrevNext}
}

func (p *playerModel) PreviousSongCmd() tea.Msg {
	s, err := p.player.PreviousSong()
	if err != nil {
		log.Fatalf("Failed to get previous song: %s", err)
	}
	return songChangedMsg{s, PrevNext}
}

var playbarTickCmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return pbTickMsg(t)
})

// tea.Model impl
func (p playerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)
	switch m := msg.(type) {
	// TODO: impl Window resize to greedily take up space with the prog bar
	case tea.KeyMsg:
		log.Printf("Got keypress: %s", m.String())
		switch {
		case key.Matches(m, playerKeys.PlayPause):
			return p, p.PlayPauseCmd
		case key.Matches(m, playerKeys.NextSong):
			return p, p.NextSongCmd
		case key.Matches(m, playerKeys.PrevSong):
			return p, p.PreviousSongCmd
		}
	case songChangedMsg:
		log.Printf("Song Updated! %v", m)
		if (p.song == pl.Song{}) {
			cmd = playbarTickCmd
		}
		// Sometimes next/previous will beat the current song request and we need to check again
		if p.song.Equals(m.song) && m.event == PrevNext {
			cmd = p.SongChangedCmd
		}
		p.song = m.song
		return p, cmd
	case pbTickMsg:
		// Approximately 1 second between ticks, we can probably calc more accurately if needed
		// TODO: Should we delegate this to the player?
		// If so, this can't be an async call and need to be managed in-band
		if p.player.IsPlaying() {
			p.song.Position += 1000
		}

		// The song is over and we should check for the new song
		if p.song.Position >= p.song.Length {
			cmd = p.SongChangedCmd
		}
		return p, tea.Sequence(cmd, playbarTickCmd)

	}
	return p, nil
}

func (p playerModel) View() string {
	// style playbar
	pbar := lip.JoinHorizontal(
		lip.Center,
		pad(formatMs(p.song.Position)),
		// using View is nice for the bounce, but this works for now
		p.playbar.ViewAs(p.songPctComplete()),
		pad(formatMs(p.song.Length)),
	)
	song := lip.NewStyle().Bold(true).Render(lip.JoinVertical(lip.Center, p.song.Name, p.song.Artist))
	return lip.JoinHorizontal(lip.Center, song, "  ", pbar)
}

func (p playerModel) Init() tea.Cmd {
	return nil
}

func NewPlayerModel() *playerModel {
	pb := progress.New(
		progress.WithoutPercentage(),
		progress.WithGradient("#73daca", "#7DCFFF"),
	)
	pb.Width = 150
	return &playerModel{
		playbar: pb,
	}
}

// Helpers
func (p *playerModel) songPctComplete() float64 {
	return float64(p.song.Position) / float64(p.song.Length)
}

func formatMs(ms int) string {
	var fmtTime string
	t := time.Duration(ms * int(time.Millisecond))
	t = t.Round(time.Second)
	h := t / time.Hour
	t -= h * time.Hour
	m := t / time.Minute
	t -= m * time.Minute
	s := t / time.Second
	// always show minutes/seconds
	fmtTime = fmt.Sprintf("%02d:%02d", m, s)
	// if hours are included, prepend hour count
	if h != 0 {
		fmtTime = fmt.Sprintf("%02d:%s", h, fmtTime)
	}
	return fmtTime
}

func pad(s string) string {
	return fmt.Sprintf(" %s ", s)
}