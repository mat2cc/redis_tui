package tui

import (
	"context"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/redis/go-redis/v9"
)

const MARGIN = 2

type ModelOptions struct {
	pretty_print_json bool
	include_types     bool
	scan_size         int64
}

type Model struct {
	redis *redis.Client

	conn_string string
	search      string
	ScanCursor  int
	Node        Node

	opts ModelOptions

	// models
	help       help.Model
	details    *Details
	search_bar *Search
	tpl        *TablePrintList
}

func InitialModel(
	redis *redis.Client,
	scanSize int64,
	pretty_print_json bool,
	include_types bool,
) *Model {
	help := help.New()
	help.ShowAll = false

	return &Model{
		redis: redis,
		details: &Details{
			key: "",
		},
		search_bar: NewSearch(),
		tpl:        NewTable(),
		help:       help,
		opts: ModelOptions{
			pretty_print_json: pretty_print_json,
			include_types:     include_types,
			scan_size:         scanSize,
		},
	}
}

type scanMsg struct {
	keys   []string
	cursor int
}

var ctx = context.Background()

func (m *Model) reset(search string) {
	m.ScanCursor = 0
	m.search = search
	m.Node = Node{}
	m.details.Reset()
}

// scan redis db using the search string
func (m *Model) Scan() tea.Cmd {
	keys, cursor, err := m.redis.Scan(ctx, uint64(m.ScanCursor), m.search, m.opts.scan_size).Result()
	if err != nil {
		log.Fatal(err)
	}
	return func() tea.Msg {
		return scanMsg{keys, int(cursor)}
	}
}

func (m *Model) Init() tea.Cmd {
	return m.Scan()
}

// set details depending on the type of the key
func (m *Model) GetDetails(node *Node) tea.Cmd {
	var rt string
	var err error
	var res RedisType

	rt = node.RedisType
	if rt == "" {
		rt, err = m.redis.Type(ctx, node.FullKey).Result()
		log.Println("rt", rt)
		if err != nil {
			log.Fatal(err)
		}
	}
	switch rt {
	case "string":
		res = GenerateStringType(m.redis, node, m.opts.pretty_print_json)
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
		return setDetailsMessage{node.FullKey, rt, res.Print(m.details.width)}
	}
}

func (m *Model) UpdateSize(width int, height int) {
	m.tpl.width = width/2 - MARGIN
	m.details.width = width/2 - MARGIN

	m.tpl.height = height - 7
	m.details.height = height - 7
}

func setWindowSize(width int, height int) tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  width,
			Height: height,
		}
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateSize(msg.Width, msg.Height)

	case scanMsg: // new scan results
		search := strings.ReplaceAll(m.search, "*", "")
		m.Node.GenNodes(msg.keys, m.redis, search, m.opts.include_types)

		m.tpl.Update(updatePL{&m.Node})
		m.ScanCursor = msg.cursor

	case setTextMessage:
		m.reset(msg.text)
		cmds := tea.Sequence(m.Scan(), m.tpl.ResetCursor)
		return m, cmds

	case tea.KeyMsg:
		// if the search bar is active, pass all messages to search
		if m.search_bar.active {
			ms, cmd := m.search_bar.Update(msg)
			m.search_bar = ms.(*Search)

			return m, cmd
		}
		switch {
		case key.Matches(msg, default_keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, default_keys.Search):
			m.search_bar.ToggleActive(true)
			return m, nil

		case key.Matches(msg, default_keys.Enter):
			node := m.tpl.GetCurrent()
			if node != nil && len(node.Children) == 0 {
				cmd := m.GetDetails(node)
				return m, cmd
			} else if node != nil {
				node.Expanded = !node.Expanded
				m.tpl.Update(updatePL{&m.Node})
			}

		case key.Matches(msg, default_keys.Scan):
			return m, m.Scan()

		case key.Matches(msg, default_keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case key.Matches(msg, default_keys.Search):
			m.search_bar.ToggleActive(true)
			return m, nil
		}
	}

	// update the other models on the screen
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
	return m, nil
}

func (m *Model) View() string {
	search := m.search_bar.View()
	print_list := m.tpl.View()

	var main string
	main = lipgloss.JoinHorizontal(lipgloss.Left,
		print_list,
		m.details.View(),
	)

	var helpmap help.KeyMap
	if m.search_bar.active {
		helpmap = search_keys
	} else {
		helpmap = default_keys
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		search,
		main,
		m.help.View(helpmap),
	)
}
