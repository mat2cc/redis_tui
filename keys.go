package main

import "github.com/charmbracelet/bubbles/key"

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k defaultKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k defaultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Enter, k.Search},
		{k.Scan},
		{k.Help, k.Quit},
	}
}

type defaultKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Search key.Binding
	Scan   key.Binding
}

var default_keys = defaultKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Search: key.NewBinding(
		key.WithKeys("/", "s"),
		key.WithHelp("s or /", "search"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "e"),
		key.WithHelp("enter/e", "enter"),
	),
	Scan: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "scan more"),
	),
}

func (k searchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Esc, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k searchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter},
		{k.Esc},
		{k.Quit},
	}
}

type searchKeyMap struct {
	Quit  key.Binding
	Esc   key.Binding
	Enter key.Binding
}

var search_keys = searchKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit search"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm search"),
	),
}

func (k redisInputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k redisInputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter},
		{k.Quit},
	}
}

type redisInputKeyMap struct {
	Quit  key.Binding
	Enter key.Binding
	Up key.Binding
	Down key.Binding
}

var redis_input_keys = redisInputKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm search"),
	),
    Up: key.NewBinding(
        key.WithKeys("up", "shift+tab"),
        key.WithHelp("↑/shift+tab", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "tab"),
        key.WithHelp("↓/tab", "move down"),
    ),
}
