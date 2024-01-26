package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Enter          key.Binding
	ProviderWindow key.Binding
	Quit           key.Binding
}

var (
	keys = keyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
		),
		ProviderWindow: key.NewBinding(
			key.WithKeys("r"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl + c", "q"),
		),
	}
)
