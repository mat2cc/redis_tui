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
	// choices  []string         // items on the to-do list
	// selected map[int]struct{} // which to-do items are selected

	pl      *PrintList
	details *Details

	redis *redis.Client

	search      string
	scan_cursor int
	node        Node

	search_bar Search
}

func initialModel() Model {
	return Model{
		// Our to-do list is a grocery list
		// choices: []string{},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		// selected: make(map[int]struct{}),
		redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		pl: &PrintList{
			List:   make([]*PrintItem, 0),
			cursor: 0,
		},
		details: &Details{
			key:  "",
			open: false,
		},
		search_bar: NewSearch(),
	}
}

type scanMsg struct {
	keys   []string
	cursor int
}

var ctx = context.Background()

func (m *Model) reset(search string) {
	m.scan_cursor = 0
	m.search = search
	m.node = Node{}
	m.pl = &PrintList{
		List:   make([]*PrintItem, 0),
		cursor: 0,
	}
}

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

func (m *Model) GetDetails(node *Node) tea.Cmd {
	key := node.FullKey
	res, err := m.redis.Get(ctx, key).Result()
	if err != nil {
		log.Fatal(err)
	}
	return func() tea.Msg {
		return setDetailsMessage{key, res}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		var model tea.Model

		model, cmd = m.search_bar.Update(msg)
		m.search_bar = model.(Search)

		model, cmd = m.pl.Update(msg)
		m.pl = model.(*PrintList)

		cmds = append(cmds, cmd)
		// return m, tea.Batch(cmds...)
	case scanMsg:
		for _, key := range msg.keys {
			split := strings.Split(key, ":")
			m.node.AddChild(split, key)
		}
		m.pl.Update(updatePL{&m.node})

		m.scan_cursor = msg.cursor
	case setTextMessage:
		m.reset(msg.text)
		return m, m.Scan()

	// Is it a key press?
	case tea.KeyMsg:
		if m.search_bar.active {
			ms, cmd := m.search_bar.Update(msg)
			m.search_bar = ms.(Search)

			return m, cmd
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

		case "enter":
      node := m.pl.List[m.pl.cursor].Node
      if len(node.Children) == 0 {
        cmd := m.GetDetails(node)
        if !m.details.open {
          m.details.open = true
        }
        return m, cmd
      }
		case "e":
			m.pl.ToggleExpand()
			m.pl.Update(updatePL{&m.node})

			// TODO: if on a leaf node, find the previous node an close expand
			// maybe look at the depth and find one less depth that at cursor
		}
	}
	res, cmd := m.pl.Update(msg)
	if a, ok := res.(*PrintList); ok {
		m.pl = a
	} else {
		return res, cmd
	}

	res, cmd = m.details.Update(msg)
	if a, ok := res.(*Details); ok {
		m.details = a
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	search := m.search_bar.View()
	print_list := m.pl.View()

	// The footer
	footer := "\nPress q to quit.\tPress d for details\n"

	var main string
	if m.details.open {
		main = lipgloss.JoinHorizontal(lipgloss.Left,
			print_list,
			m.details.View(),
		)
	} else {
		main = print_list
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		search,
		main,
		footer,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
