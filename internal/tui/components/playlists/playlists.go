package playlists

import (
	"log"
	"tmix/internal/providers"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	p         providers.MusicProvider
	playlists list.Model
}

// Playlist list items
type playlistListItem struct {
	providers.Playlist
}

func (pl playlistListItem) Title() string {
	return pl.Name
}

func (pl playlistListItem) Description() string {
	return pl.Playlist.Description
}

func (pl playlistListItem) FilterValue() string {
	return pl.Name
}

// Playlist cmd
type playlistsLoadedMsg struct{ pl []providers.Playlist }

func (m *Model) LoadPlaylistsCmd() tea.Msg {
	return playlistsLoadedMsg{m.p.FetchPlaylists()}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.playlists, cmd = m.playlists.Update(msg)
	case playlistsLoadedMsg:
		var playlists []list.Item
		for _, pl := range msg.pl {
			playlists = append(playlists, playlistListItem{
				pl,
			})
		}
		m.playlists = list.New(playlists, list.NewDefaultDelegate(), 20, 10)
		m.playlists.Title = "Playlists"
		log.Printf("%d Playlists retreived", len(m.playlists.Items()))
	}
	return m, cmd
}

func (m Model) View() string {
	return m.playlists.View()
}

func New(p providers.MusicProvider) *Model {
	log.Printf("Creating new playlist model from with provider %p", p)
	return &Model{
		p,
		list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
}
