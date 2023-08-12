package main

import "github.com/charmbracelet/bubbles/key"

type helpKeyMap map[string]key.Binding

func (k helpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("Ctrl+C"), key.WithHelp("Ctrl+C", "quit")),
		key.NewBinding(key.WithKeys("Tab"), key.WithHelp("Tab", "send request")),
		key.NewBinding(key.WithKeys("Ctrl+Up"), key.WithHelp("Ctrl+Up", "scroll result ↑")),
		key.NewBinding(key.WithKeys("Ctrl+Down"), key.WithHelp("Ctrl+Down", "scroll result ↓")),
		key.NewBinding(key.WithKeys("Alt+A"), key.WithHelp("Alt+A", "copy response")),
		key.NewBinding(key.WithKeys("Alt+S"), key.WithHelp("Alt+S", "copy headers")),
		key.NewBinding(key.WithKeys("Alt+D"), key.WithHelp("Alt+D", "copy selected header")),
	}
}

func (k helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
