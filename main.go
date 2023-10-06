package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/redis/go-redis/v9"
)

type Model struct {
	choices  []string // items on the to-do list
	cursor   int
	selected map[int]struct{} // which to-do items are selected
    print_list []*PrintItem

	redis *redis.Client

	search      string
	scan_cursor int
	node       Node
}

func (m *Model) Print() string {
	str := ""
	for i, item := range m.print_list {
        if i == m.cursor {
            str += ">"
        } else {
            str += " "
        }
        str += item.Print()
	}

	return str
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
        m.print_list = GeneratePrintList(&m.node, 0)
		m.scan_cursor = msg.cursor

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "m":
			return m, m.Scan()

        case "e":
            n := m.print_list[m.cursor].Node
            n.expanded = !n.expanded
            m.print_list = GeneratePrintList(&m.node, 0)

            // TODO: if on a leaf node, find the previous node an close expand
            // maybe look at the depth and find one less depth that at cursor

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.print_list)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	// for i, choice := range m.choices {

	// 	// Is the cursor pointing at this choice?
	// 	cursor := " " // no cursor
	// 	if m.cursor == i {
	// 		cursor = ">" // cursor!
	// 	}

	// 	// Is this choice selected?
	// 	checked := " " // not selected
	// 	if _, ok := m.selected[i]; ok {
	// 		checked = "x" // selected!
	// 	}

	// 	// Render the row
	// 	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	// }

    s += m.Print()

	// The footer
	s += "\nPress q to quit.\n"
	// s += fmt.Sprintf("%v\n\n", m.node)
    // s += fmt.Sprintf(m.node.Print(2))

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
