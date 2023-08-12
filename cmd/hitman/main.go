package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ramitmittal/hitman/internal/httpclient"
)

var (
	// injected at build time
	GitSHA string
	GitTag string
)

type model struct {
	ready        bool
	windowWidth  int
	windowHeight int

	// plain text for the title bar
	titlePlainText string

	// title bar; rendered at the top of the page
	titleComponent string

	// error bar; rendered below the title bar if last operation errored;
	errComponent string

	// help text; rendered at the bottom
	helpComponent string

	// result of an HTTP calls stored as r1, r2, r3, ...rn, \n, R1, R2, R3, ...Rn, \n, RB]
	// where r1, ...rn are request headers, R1, ...Rn are response headers, and RB is response body
	rawResult []string

	// the index of rawResult that contains the line selected in the viewport
	viewportSelectedLineIndex int

	// viewport for results of http calls; rendered below error component
	viewport viewport.Model

	// text area for user input; rendered below result viewport
	textarea textarea.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if !m.ready {
		return "Starting..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n\n%s",
		m.titleComponent,
		m.errComponent,
		m.viewport.View(),
		m.textarea.View(),
		m.helpComponent,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var stopPropogation bool

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.resetTitle()
		m.initViewport()
		m.initTextarea()
		m.initHelp()
		m.ready = true

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			saveText(m.textarea.Value())
			return m, tea.Quit

		case tea.KeyTab:
			if result := httpclient.Hit(m.textarea.Value()); result.Err != nil {
				m.setError(result.Err)
				m.viewport.SetContent("")
			} else {
				if m.errComponent != "" {
					m.unsetError()
				}
				m.setResult(result)
			}
			stopPropogation = true
		case tea.KeyCtrlDown:
			m.scrollDown()
			stopPropogation = true
		case tea.KeyCtrlUp:
			m.scrollUp()
			stopPropogation = true
		case tea.KeyRunes:
			if msg.Alt {
				switch string(msg.Runes) {
				case "a":
					m.copyResult()
				case "s":
					m.copyHeaders()
				case "d":
					m.copyHighlight()
				}
				stopPropogation = true
			}
		}
	}

	var cmds []tea.Cmd
	if !stopPropogation {
		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)

		var taCmd tea.Cmd
		m.textarea, taCmd = m.textarea.Update(msg)
		cmds = append(cmds, taCmd)
	}
	return m, tea.Batch(cmds...)
}

// Initialize the title bar
func (m *model) resetTitle() {
	m.titleComponent = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("2")).
		Width(m.windowWidth).
		Render(m.titlePlainText)
}

// Set title bar with RED background
func (m *model) errorTitle() {
	m.titleComponent = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("9")).
		Width(m.windowWidth).
		Render(m.titlePlainText)
}

// Set value for error component
func (m *model) setError(err error) {
	m.errComponent = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(err.Error())
	m.errorTitle()
	m.viewport.SetContent("")
	m.viewportSelectedLineIndex = 0
}

// Unset value for error component
func (m *model) unsetError() {
	m.errComponent = ""
	m.resetTitle()
}

// Initialize the viewport component
func (m *model) initViewport() {
	m.viewport = viewport.New(m.windowWidth, m.windowHeight*78/100)
	m.viewport.KeyMap = viewport.KeyMap{}
	m.viewport.SetContent("")
}

// Initialize the textarea component
func (m *model) initTextarea() {
	m.textarea = textarea.New()
	m.textarea.SetWidth(m.windowWidth)
	m.textarea.SetHeight(m.windowHeight * 15 / 100)
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
}

// Initialize the help component
func (m *model) initHelp() {
	h := help.New()
	h.Width = m.windowWidth
	m.helpComponent = h.View(helpKeyMap{})
}

// Attempts to copy last hit's result to clipboard; populates error component on failure
func (m *model) copyResult() {
	if len(m.rawResult) == 0 {
		m.setError(errors.New("no result to copy"))
	} else if err := copyText(strings.Join(m.rawResult, "\n")); err != nil {
		m.setError(err)
	} else {
		m.unsetError()
	}
}

// Attempts to copy last hit's response headers to clipboard; populates error component on failure
func (m *model) copyHeaders() {
	if len(m.rawResult) == 0 {
		m.setError(errors.New("no headers to copy"))
	} else if err := copyText(strings.Join(m.rawResult[:len(m.rawResult)-2], "\n")); err != nil {
		m.setError(err)
	} else {
		m.unsetError()
	}
}

// Attempts to copy the viewport's highlighted text to clipboard; populates error component on failure
func (m *model) copyHighlight() {
	if err := copyText(m.rawResult[m.viewportSelectedLineIndex]); err != nil {
		m.setError(err)
	} else {
		m.unsetError()
	}
}

// Transform httpclient.HitResult into []string and update model
func (m *model) setResult(result *httpclient.HitResult) {
	rawResult := make([]string, 0, len(result.RequestHeaders)+len(result.ResponseHeaders)+3)

	rawResult = append(rawResult, result.RequestHeaders...)
	rawResult = append(rawResult, "\n")
	rawResult = append(rawResult, result.ResponseHeaders...)
	rawResult = append(rawResult, "\n")
	rawResult = append(rawResult, result.ResponseBody)

	if m.viewportSelectedLineIndex > len(rawResult) {
		m.viewportSelectedLineIndex = 0
	}
	m.rawResult = rawResult
	m.updateFormattedResult()
}

// Convert raw http result []string into formatted text for viewport
func (m *model) updateFormattedResult() {
	highlightedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("#FFFFFF"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	vhli := m.viewportSelectedLineIndex
	rawResultLength := len(m.rawResult)

	var formattedResult strings.Builder

	for idx, head := range m.rawResult[:rawResultLength-2] {
		if idx == vhli {
			formattedResult.WriteString(highlightedStyle.Render(head))
		} else {
			formattedResult.WriteString(headerStyle.Render(head))
		}
		if head != "\n" {
			formattedResult.WriteRune('\n')
		}
	}
	formattedResult.WriteString(m.rawResult[rawResultLength-2])
	formattedResult.WriteString(m.rawResult[rawResultLength-1])

	m.viewport.SetContent(formattedResult.String())
}

func (m *model) scrollDown() {
	if m.viewportSelectedLineIndex < len(m.rawResult)-3 {
		m.viewportSelectedLineIndex += 1
		m.updateFormattedResult()
	}
	m.viewport.LineDown(1)
}

func (m *model) scrollUp() {
	if m.viewportSelectedLineIndex > 0 {
		m.viewportSelectedLineIndex -= 1
		m.updateFormattedResult()
	}
	m.viewport.LineUp(1)
}

// Returns plain text for the title bar component
func generateTitlePlainText() string {
	var version string

	if GitTag != "" {
		version = GitTag
	} else if GitSHA != "" {
		version = GitSHA
	} else {
		version = "?"
	}

	return "Hitman HTTP Client v" + version
}

func main() {
	m := model{
		titlePlainText: generateTitlePlainText(),
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
