package cli

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	NONE mode = iota
	ADD
	EDIT
)

type textModel struct {
	mainModel *mainModel
	input     textinput.Model
	mode      mode
	help      help.Model
	keys      *textKeyMap
	focus     bool
}

func initText(main *mainModel) *textModel {
	t := textinput.New()
	t.CharLimit = 255
	t.Width = 100

	return &textModel{
		mainModel: main,
		input:     t,
		mode:      NONE,
		help:      help.New(),
		keys:      newTextInputKeys(),
		focus:     false,
	}
}

func (m textModel) Init() tea.Cmd { return nil }

func (m *textModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.Focussed() {
			return nil
		}

	case tea.WindowSizeMsg:
		h, _ := appStyle.GetFrameSize()
		m.input.Width = msg.Width - h
	}

	hasTyped := len(m.input.Value()) != 0
	m.keys.confirm.SetEnabled(hasTyped)

	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m textModel) View() string {
	v := lipgloss.JoinVertical(lipgloss.Left, m.input.View(), m.help.View(m.keys))
	return textStyle.Render(v)
}

func (m textModel) Focussed() bool {
	return m.focus
}

// Focus the text input and select the mode
func (m *textModel) Focus(mode mode, name string) {
	m.focus = true
	m.mode = mode
	m.input.Placeholder = "To-do name"
	if mode == EDIT {
		m.input.SetValue(name)
	}
	m.input.Focus()
}

func (m *textModel) Blur() {
	m.focus = false
	m.mode = NONE
	m.input.Reset()
	m.input.Blur()
}

func (m *textModel) Mode() mode { return m.mode }
