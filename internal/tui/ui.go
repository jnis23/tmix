package tui

import (
	"tmix/internal/providers"
	"tmix/internal/tui/components/playlists"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	Normal Mode = iota
	Mini
)

type model struct {
	width, height int
	mode          Mode
	// Limit to one provider for now
	currentProvider providers.MusicProvider
	currentContext  tea.Model
	providers       providersModel
	player          playerModel
	playlists       playlists.Model
}

var (
	widgetStyle = lip.NewStyle().
			BorderStyle(lip.RoundedBorder()).
			BorderForeground(lip.Color("63")).
			Align(lip.Center)
	focusedWidget = lip.NewStyle().
			BorderStyle(lip.DoubleBorder()).
			Inherit(widgetStyle)
)

// tea.Model impl
func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: Resize components
		// player should always be shown and should take up the entire top of the screen
		// playlists should be on left side similar to Spotify UI
		// song results should take up the remaining center view
		m.height, m.width = msg.Height, msg.Width
		m.providers.providers.SetHeight(m.height)
		pv, cmd := m.providers.Update(msg)
		m.providers = pv.(providersModel)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.ProviderWindow):
			m.currentContext = m.providers
		}
		if m.currentProvider == nil {
			_, cmd := m.providers.Update(msg)
			cmds = append(cmds, cmd)
		}
	case loginCompleteMsg:
		m.currentProvider = msg
		m.player.player = m.currentProvider.Player()
		if m.mode == Normal {
			m.playlists = *playlists.New(m.currentProvider)
			cmd = m.playlists.LoadPlaylistsCmd
		}
		// Now that we are logged in, we can prepopulate current song and playlist data
		cmds = append(cmds, m.player.SongChangedCmd, cmd)
	}

	if m.currentProvider != nil {
		pl, cmd := m.player.Update(msg)
		m.player = pl.(playerModel)
		cmds = append(cmds, cmd)

		if m.mode == Normal {
			playlist, cmd := m.playlists.Update(msg)
			if playlist, ok := playlist.(playlists.Model); ok {
				m.playlists = playlist
			}
			cmds = append(cmds, cmd)
		}

	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// style all widgets with base style
	// update only current context with focusedStyle
	return m.MiniMode()

	// player := widgetStyle.Copy().Width(m.width - 10).Render(m.player.View())
	// providers := widgetStyle.Render(m.providers.View())
	// playlist := widgetStyle.Copy().Width(50).Render(m.playlists.View())
	// switch m.currentContext.(type) {
	// case providersModel:
	// 	providers = focusedWidget.Render(m.providers.View())
	// }
	// core := lip.JoinHorizontal(lip.Left, providers, playlist)
	// return lip.JoinVertical(lip.Left, player, core)
}

// MiniMode will strictly provide login to a provider and access to the player.
func (m model) MiniMode() string {
	if m.currentProvider == nil {
		return widgetStyle.Render(m.providers.View())
	} else {
		return widgetStyle.Copy().Width(m.width - 10).Render(m.player.View())
	}
}

type Config struct {
	Providers *providers.ProviderConfig `toml:"providers"`
}

func New(config *Config) (*model, error) {
	// get a list of components that can be focused
	return &model{
		width:     0,
		height:    0,
		providers: *NewProviders(config.Providers),
		player:    *NewPlayerModel(),
	}, nil
}
