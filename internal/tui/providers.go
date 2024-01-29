package tui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"tmix/internal/config"
	pro "tmix/internal/providers"
)

type providersModel struct {
	currProvider pro.MusicProvider
	providers    list.Model
	logWindow    logModel
}

type providerListItem struct {
	pro.MusicProvider
}

func (p providerListItem) Title() string { return p.Name() }
func (p providerListItem) Description() string {
	if p.LoggedIn() {
		return "Logged In"
	}
	return "Log In"
}
func (p providerListItem) FilterValue() string { return p.Name() }

func (m providersModel) Init() tea.Cmd {
	return nil
}

type tickMsg time.Time
type loginCompleteMsg pro.MusicProvider

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *providersModel) loginCmd(pvi *providerListItem) tea.Cmd {
	return func() tea.Msg {
		pvi.Login()
		log.Printf("Login completed with client: %s", pvi.Name())
		log.Printf("LoginMsg provider %p", pvi)
		m.currProvider = pvi
		return loginCompleteMsg(pvi)
	}
}

func (m providersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			s := m.providers.SelectedItem().(providerListItem)
			if s.LoggedIn() {
				return m, nil
			}
			log.Printf("Logging into %s", s.Name())
			return m, m.loginCmd(&s)
		}

	}

	var cmd tea.Cmd
	// forward to list
	m.providers, cmd = m.providers.Update(msg)
	cmds = append(cmds, cmd)
	lw, cmd := m.logWindow.Update(msg)
	m.logWindow = lw.(logModel)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m providersModel) View() string {
	return m.providers.View()
}

func NewProviders(config *config.ProviderConfig) *providersModel {
	items := []list.Item{}
	for _, p := range pro.LoadProviders(config) {
		pv := providerListItem{p}
		items = append(items, pv)
	}
	l := list.New(items, list.NewDefaultDelegate(), 20, 10)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	return &providersModel{
		providers: l,
		logWindow: *newLogModel(),
	}
}
