package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mat2cc/redis_tui/redis_type"
	"github.com/redis/go-redis/v9"
)

const MARGIN = 2 

type Model struct {
	redis *redis.Client

	search      string
	scan_cursor int
	node        Node

	// models
	details    *Details
	search_bar *Search
	tpl        *TablePrintList
}

func initialModel() Model {
	return Model{
		redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		details: &Details{
			key:  "",
			open: true,
		},
		search_bar: NewSearch(),
		tpl:        NewTable(),
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
	m.details = &Details{
		key:   "",
		open:  m.details.open,
		width: m.details.width,
	}
}

func (m *Model) Scan() tea.Cmd {
	keys, cursor, err := m.redis.Scan(ctx, uint64(m.scan_cursor), m.search, 1000).Result()
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
	rt, err := m.redis.Type(ctx, node.FullKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	var res redis_type.RedisType
	switch rt {
	case "string":
		res = GenerateStringType(m.redis, node)
	case "list":
		res = GenerateListType(m.redis, node)
	case "set":
		res = GenerateSetType(m.redis, node)
	case "zset":
		res = GenerateZSetType(m.redis, node)
	case "hash":
		res = GenerateHashType(m.redis, node)
	case "stream":
		res = GenerateStreamType(m.redis, node)
	}

	return func() tea.Msg {
		if res == nil {
			return setDetailsMessage{node.FullKey, "", "Type Not implemented"}
		}
		return setDetailsMessage{node.FullKey, node.RedisType, res.Print(m.details.width)}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
    m.tpl.width = msg.Width / 2 - MARGIN
    m.details.width = msg.Width / 2 - MARGIN

	case scanMsg:
		for _, key := range msg.keys {
			split := strings.Split(key, ":")
			m.node.AddChild(split, key, m.redis)
		}
		m.tpl.Update(updatePL{&m.node})

		m.scan_cursor = msg.cursor
	case setTextMessage:
		m.reset(msg.text)
		return m, m.Scan()

	// Is it a key press?
	case tea.KeyMsg:
		if m.search_bar.active {
			ms, cmd := m.search_bar.Update(msg)
			m.search_bar = ms.(*Search)

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
			m.search_bar.ToggleActive(true)
      return m, nil

		case "enter":
			node := m.tpl.GetCurrent()
			if node != nil && len(node.Children) == 0 {
				cmd := m.GetDetails(node)
				if !m.details.open {
					m.details.open = true
				}
				return m, cmd
			}
		case "e":
			m.tpl.ToggleExpand()
			m.tpl.Update(updatePL{&m.node})

			// TODO: if on a leaf node, find the previous node an close expand
			// maybe look at the depth and find one less depth that at cursor
		}
	}

	res, cmd := m.tpl.Update(msg)
	if a, ok := res.(*TablePrintList); ok {
		m.tpl = a
	} else {
		return res, cmd
	}

	res, cmd = m.search_bar.Update(msg)
	if a, ok := res.(*Search); ok {
		m.search_bar = a
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
	print_list := m.tpl.View()

	// The footer
	footer := "Press q to quit.\tPress d for details"

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
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
