package cli

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zGraund/TodoCli/internal/models"
)

type action int

const (
	REFRESH action = iota // Refresh the todo list
	INSERT                // Insert new todo
	RENAME                // Rename a todo
	REMOVE                // Remove a todo
	UPDATE                // Toggle status
)

type todoMsg struct {
	action action
	todos  []models.Todo
}

type listModel struct {
	main   *mainModel
	todos  list.Model
	date   time.Time
	offset int
	keys   *listKeyMap
	focus  bool
}

func (m listModel) Init() tea.Cmd { return m.updateTodos }

func (m *listModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.Focussed() {
			return nil
		}
		if m.todos.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.prevDay):
			m.offset--
			return m.updateTodos
		case key.Matches(msg, m.keys.nextDay):
			m.offset++
			return m.updateTodos
		case key.Matches(msg, m.keys.markDon):
			return m.markAsDone
		case key.Matches(msg, m.keys.remTodo):
			return m.removeTodo
		}

	case todoMsg:
		switch msg.action {
		case REFRESH:
			items := make([]list.Item, len(msg.todos))
			for i, todo := range msg.todos {
				items[i] = todo
			}
			cmds = append(cmds, m.todos.SetItems(items))
		case INSERT:
			index := len(m.todos.Items())
			cmds = append(cmds, m.todos.InsertItem(index, msg.todos[0]), m.todos.NewStatusMessage("New To-do added"))
		case RENAME:
			cmds = append(cmds, m.todos.SetItem(m.todos.Index(), msg.todos[0]), m.todos.NewStatusMessage("Todo name updated"))
		case REMOVE:
			m.todos.RemoveItem(m.todos.Index())
			cmds = append(cmds, m.todos.NewStatusMessage("To-do removed"))
		case UPDATE:
			cmds = append(cmds, m.todos.SetItem(m.todos.Index(), msg.todos[0]), m.todos.NewStatusMessage(msg.todos[0].Status()))
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.todos.SetSize(msg.Width-h, msg.Height-v)
	}

	m.todos, cmd = m.todos.Update(msg)
	cmds = append(cmds, cmd)
	m.updateKeyMap()
	return tea.Batch(cmds...)
}

func (m listModel) View() string {
	dayOffset := time.Hour * 24 * time.Duration(m.offset)
	m.todos.Title = m.date.Add(dayOffset).Format("Monday 02/01/2006")
	if m.offset == 0 {
		m.todos.Title = "Today " + m.todos.Title
	}
	return m.todos.View()
}

func initList(main *mainModel) *listModel {
	date := time.Now()
	keys := newListKeyMap()

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemTitle
	delegate.Styles.SelectedDesc = selectedItemDesc
	delegate.Styles.FilterMatch = filterMatch
	delegate.FullHelpFunc = keys.FullHelp
	delegate.ShortHelpFunc = keys.ShortHelp

	list := list.New([]list.Item{}, delegate, 0, 0)
	list.StatusMessageLifetime = time.Second * 2
	list.Styles.NoItems.PaddingLeft(2)
	list.SetStatusBarItemName("Todo", "Todos")

	return &listModel{
		main:   main,
		date:   date,
		offset: 0,
		todos:  list,
		keys:   keys,
		focus:  true,
	}
}

func (m listModel) updateTodos() tea.Msg {
	t, err := models.GetByDay(m.date, m.offset)
	if err != nil {
		panic("Error while querying db for todos")
	}
	return todoMsg{action: REFRESH, todos: t}
}

func (m *listModel) addTodo(name string) tea.Cmd {
	return func() tea.Msg {
		todo, err := models.Create(name, m.offset)
		if err != nil {
			panic(fmt.Sprintf("Error while inserting new todo:\n\n%v", err))
		}
		return todoMsg{action: INSERT, todos: []models.Todo{todo}}
	}
}

func (m listModel) removeTodo() tea.Msg {
	todo, ok := m.todos.SelectedItem().(models.Todo)
	if !ok {
		return nil
	}
	if err := todo.Delete(); err != nil {
		panic(fmt.Sprintf("Error while deleting todo:\n\n%v", err))
	}
	return todoMsg{action: REMOVE}
}

func (m *listModel) editTodo(name string) tea.Cmd {
	return func() tea.Msg {
		todo, ok := m.todos.SelectedItem().(models.Todo)
		if !ok {
			return nil
		}
		if err := todo.SetName(name); err != nil {
			panic(fmt.Sprintf("Error while changing todo name:\n\n%v", err))
		}
		return todoMsg{action: RENAME, todos: []models.Todo{todo}}
	}
}

func (m listModel) markAsDone() tea.Msg {
	todo, ok := m.todos.SelectedItem().(models.Todo)
	if !ok {
		return nil
	}
	if e := todo.SetStatus(); e != nil {
		panic("Error while changing Todo status")
	}
	return todoMsg{action: UPDATE, todos: []models.Todo{todo}}
}

func (m *listModel) updateKeyMap() {
	hasTodos := len(m.todos.Items()) != 0
	m.todos.SetShowStatusBar(hasTodos)
	m.keys.markDon.SetEnabled(hasTodos)
	m.keys.remTodo.SetEnabled(hasTodos)
	m.keys.chgName.SetEnabled(hasTodos)
}

func (m listModel) Focussed() bool {
	return m.focus
}

func (m *listModel) Focus() {
	m.focus = true

	// Reverse the changes from the Blur() function
	m.todos.SetShowHelp(true)
	m.todos.SetHeight(m.todos.Height() + 2)
}

func (m *listModel) Blur() {
	m.focus = false

	// For formatting purposes disable the list help menu
	// and decrease the list size to fit nice in the screen
	m.todos.SetShowHelp(false)
	m.todos.SetHeight(m.todos.Height() - 2)
}
