package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/redis/go-redis/v9"
)

type Model struct {
	choices  []string         // items on the to-do list
	selected map[int]struct{} // which to-do items are selected

	pl *PrintList

	redis *redis.Client

	search      string
	scan_cursor int
	node        Node

    search_bar Search
}

func initialModel() Model {
	return Model{
		// Our to-do list is a grocery list
		choices: []string{},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
		redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		node: Node{
			Value: "*", Children: make([]*Node, 0),
		},
		pl: &PrintList{
			List:   make([]*PrintItem, 0),
			cursor: 0,
		},
        search_bar: NewSearch(),
	}
}

type scanMsg struct {
	keys   []string
	cursor int
}

var ctx = context.Background()

func (m *Model) Scan() tea.Cmd {
	keys, cursor, err := m.redis.Scan(ctx, uint64(m.scan_cursor), m.search, 10).Result()
	if err != nil {
		log.Fatal(err)
	}
	return func() tea.Msg {
		return scanMsg{keys, int(cursor)}
	}
}

func (m Model) Init() tea.Cmd {
	return m.Scan()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case scanMsg:
		for _, key := range msg.keys {
			split := strings.Split(key, ":")
			m.node.AddChild(split[1:])
		}
		m.choices = append(m.choices, msg.keys...)
		// m.print_list = GeneratePrintList(&m.node, 0)
		m.pl.Update(updatePL{&m.node})

		m.scan_cursor = msg.cursor

	// Is it a key press?
	case tea.KeyMsg:
        if m.search_bar.active {
            return m.search_bar.Update(msg)
        }

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "m":
			return m, m.Scan()

        case "s":
            m.search_bar.active = true
            m.search_bar.input.Focus()

		case "e":
			m.pl.ToggleExpand()
			m.pl.Update(updatePL{&m.node})

			// TODO: if on a leaf node, find the previous node an close expand
			// maybe look at the depth and find one less depth that at cursor

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			// case "enter", " ":
			// 	_, ok := m.selected[m.cursor]
			// 	if ok {
			// 		delete(m.selected, m.cursor)
			// 	} else {
			// 		m.selected[m.cursor] = struct{}{}
			// 	}
			// }
		}
	}
	res, cmd := m.pl.Update(msg)
	if a, ok := res.(*PrintList); ok {
		m.pl = a
	} else {
		return res, cmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	main := m.pl.View()

	// The footer
	main += "\nPress q to quit.\n"

	search := m.search_bar.View()

	// Send the UI for rendering
	return lipgloss.JoinVertical(lipgloss.Top,
		search,
		main)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
