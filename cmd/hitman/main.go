package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ramitmittal/hitman/internal/httpclient"
)

var GitSHA string

type helpKeyMap map[string]key.Binding

func (k helpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("Tab"), key.WithHelp("Tab", "send request")),
		key.NewBinding(key.WithKeys("Ctrl+S"), key.WithHelp("Ctrl+S", "copy response")),
		key.NewBinding(key.WithKeys("Ctrl+C"), key.WithHelp("Ctrl+C", "quit")),
	}
}

func (k helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}

var hkm = helpKeyMap{}

type model struct {
	ready bool

	// static title line
	title string

	// text for the viewport (formatted req+res headers)
	messages string

	// result from the last HTTP call
	hResult httpclient.HitResult

	textarea textarea.Model
	viewport viewport.Model
	help     help.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if !m.ready {
		return "\n Starting..."
	}
	return fmt.Sprintf(
		"%s\n%s\n\n%s\n%s",
		m.title,
		m.viewport.View(),
		m.textarea.View(),
		m.help.View(hkm),
	)
}

func (m *model) prepMessagesForViewport() {
	if m.hResult.Err != nil {
		m.messages = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.hResult.Err.Error())
		return
	}

	requestHeaders := m.hResult.RequestHeaders()
	responseHeaders := m.hResult.ResponseHeaders()

	var messages strings.Builder
	messages.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("#FFFFFF")).Render(requestHeaders[0]))
	messages.WriteString("\n")
	messages.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Render(strings.Join(requestHeaders[1:], "\n")))
	messages.WriteString("\n\n")
	messages.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("#FFFFFF")).Render(responseHeaders[0]))
	messages.WriteString("\n")
	messages.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Render(strings.Join(responseHeaders[1:], "\n")))
	m.messages = messages.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updateViewport bool
	var updateTextarea bool

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.title = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("2")).
			Width(msg.Width).
			Render(fmt.Sprintf("Hitman Alpha. Build Version: %s", GitSHA))

		m.viewport = viewport.New(msg.Width, msg.Height*75/100)
		m.viewport.SetContent(m.messages)

		m.textarea = textarea.New()
		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(msg.Height * 13 / 100)
		m.textarea.Prompt = "┃ "
		m.textarea.FocusedStyle.CursorLine = lipgloss.NewStyle()
		m.textarea.ShowLineNumbers = false

		if !m.ready {
			m.textarea.SetValue(loadText())
		}
		for i := 0; i < m.textarea.LineCount(); i++ {
			m.textarea.CursorUp()
		}
		m.textarea.CursorEnd()
		m.textarea.Focus()

		m.help = help.New()
		m.help.Width = msg.Width
		m.ready = true

		updateViewport = true
		updateTextarea = true

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			// TODO: handle error
			copyText(m.hResult.String())
		case tea.KeyCtrlC, tea.KeyEsc:
			saveText(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyTab:
			// TODO: this should be async
			m.hResult = httpclient.Hit(m.textarea.Value())
			m.prepMessagesForViewport()
			m.viewport.SetContent(m.messages)
			m.viewport.GotoTop()
		default:
			updateViewport = true
			updateTextarea = true
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseWheelDown || msg.Type == tea.MouseWheelUp {
			updateViewport = true
		}
	}

	var cmds []tea.Cmd
	if updateViewport {
		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}
	if updateTextarea {
		var taCmd tea.Cmd
		m.textarea, taCmd = m.textarea.Update(msg)
		cmds = append(cmds, taCmd)
	}

	return m, tea.Batch(cmds...)
}

func main() {
	p := tea.NewProgram(model{}, tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
