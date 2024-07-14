package cli

import (
	"github.com/charmbracelet/bubbles/key"
)

type listKeyMap struct {
	prevDay key.Binding
	nextDay key.Binding
	addTodo key.Binding
	// TODO: global todo: addGlob key.Binding
	remTodo key.Binding
	chgName key.Binding
	markDon key.Binding
}

func (k *listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.markDon,
	}
}

func (k *listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.nextDay,
			k.prevDay,
		},
		{
			k.addTodo,
			// k.addGlob,
			k.remTodo,
			k.chgName,
			k.markDon,
		},
	}
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		nextDay: key.NewBinding(
			key.WithKeys("shift+right", "L"),
			key.WithHelp("shift+→", "next day"),
		),
		prevDay: key.NewBinding(
			key.WithKeys("shift+left", "H"),
			key.WithHelp("shift+←", "previous day"),
		),
		addTodo: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add todo"),
		),
		// addGlob: key.NewBinding(
		// 	key.WithKeys("A"),
		// 	key.WithHelp("A", "add global todo"),
		// ),
		remTodo: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete todo"),
		),
		chgName: key.NewBinding(
			key.WithKeys("E"),
			key.WithHelp("E", "edit todo"),
		),
		markDon: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("↵/space", "mark as done"),
		),
	}
}

type textKeyMap struct {
	confirm key.Binding
	cancel  key.Binding
}

func (k *textKeyMap) FullHelp() [][]key.Binding { return nil }

func (k *textKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.confirm,
		k.cancel,
	}
}

func newTextInputKeys() *textKeyMap {
	return &textKeyMap{
		confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "confirm"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}
