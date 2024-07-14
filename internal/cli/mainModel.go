package cli

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zGraund/TodoCli/internal/models"
)

type mainModel struct {
	list *listModel
	text *textModel
}

func InitMainModel() mainModel {
	m := mainModel{}
	m.list = initList(&m)
	m.text = initText(&m)
	return m
}

func (m mainModel) Init() tea.Cmd { return m.list.Init() }

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		switch {
		case m.list.Focussed() && m.list.todos.FilterState() != list.Filtering:
			switch {
			case key.Matches(msg, m.list.keys.addTodo):
				m.list.Blur()
				m.text.Focus(ADD, "")
				return m, textinput.Blink
			case key.Matches(msg, m.list.keys.chgName):
				todo := m.list.todos.SelectedItem().(models.Todo)
				m.list.Blur()
				m.text.Focus(EDIT, todo.Name)
				return m, textinput.Blink
			}
		case m.text.Focussed():
			if key.Matches(msg, m.text.keys.cancel) {
				m.list.Focus()
				m.text.Blur()
				return m, nil
			}
			if key.Matches(msg, m.text.keys.confirm) {
				name := m.text.input.Value()
				switch m.text.Mode() {
				case ADD:
					cmd = m.list.addTodo(name)
				case EDIT:
					cmd = m.list.editTodo(name)
				}
				m.list.Focus()
				m.text.Blur()
				return m, cmd
			}
		}
	}

	cmds = append(cmds, m.list.Update(msg), m.text.Update(msg))
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	if m.text.Focussed() {
		s = lipgloss.JoinVertical(lipgloss.Left, m.list.View(), m.text.View())
	} else {
		s = m.list.View()
	}
	return appStyle.Render(s)
}
