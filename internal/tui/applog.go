package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type logModel struct {
	w  applog
	vp viewport.Model
}

type applog struct{}

var appLog = ""

func (w *applog) Write(data []byte) (n int, err error) {
	dat := string(data)
	appLog += "\n" + dat
	return len(data), nil
}

func (m logModel) Init() tea.Cmd {
	return nil
}

func (m logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.vp.SetContent(wordwrap.String(appLog, lipgloss.Width(m.vp.View())))
	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m logModel) View() string {
	return widgetStyle.Render(m.vp.View())
}

func newLogModel() *logModel {
	vp := viewport.New(100, 50)
	vp.HighPerformanceRendering = false
	lm := logModel{
		w:  applog{},
		vp: vp,
	}
	//log.SetOutput(&lm.w)
	log.Printf("New logger set")
	return &lm
}
