package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/redis/go-redis/v9"
)

type RedisInput struct {
	fields      []textinput.Model
	focusCursor int

	width  int
	height int
	error  string
}

func NewRedisInput() *RedisInput {
	connection := textinput.New()
	connection.Placeholder = "redis://localhost:6379"
	connection.Prompt = "Connection: "
	connection.Focus()

	username := textinput.New()
	username.Prompt = "Username (optional): "

	password := textinput.New()
	password.Prompt = "Password (optional): "

	return &RedisInput{
		fields: []textinput.Model{connection, username, password},
	}
}

func createRedisClient(conn string) (*redis.Client, error) {
	if conn == "" {
		conn = "localhost:6379"
	} else {
		conn = strings.TrimPrefix(conn, "redis://")
	}
	redis := redis.NewClient(&redis.Options{
		Addr:     conn,
		Username: "",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return redis, nil
}

func (i *RedisInput) Init() tea.Cmd {
	return textinput.Blink
}

func (i *RedisInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.width = msg.Width
		i.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, redis_input_keys.Quit):
			return i, tea.Quit
		case key.Matches(msg, redis_input_keys.Enter):
			redis, err := createRedisClient(i.fields[0].Value())
			if err != nil {
				i.error = err.Error()
				return i, nil
			}

			model := initialModel(redis)
			cmds := tea.Batch(setWindowSize(i.width, i.height), model.Scan())
			return model, cmds
		case key.Matches(msg, redis_input_keys.Up):
			if i.focusCursor > 0 {
				i.focusCursor--
			}
			i.UpdateInputs()
		case key.Matches(msg, redis_input_keys.Down):
			if i.focusCursor < len(i.fields)-1 {
				i.focusCursor++
			}
			i.UpdateInputs()
		}
	}

	cmd = i.UpdateFields(msg)

	return i, cmd
}

func (ri *RedisInput) UpdateInputs() tea.Cmd {
	cmds := make([]tea.Cmd, len(ri.fields))

	for i := range ri.fields {
		if i == ri.focusCursor {
			cmds[i] = ri.fields[i].Focus()
			continue
		}
		ri.fields[i].Blur()
	}
	return tea.Batch(cmds...)
}

func (ri *RedisInput) UpdateFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(ri.fields))

	for i := range ri.fields {
		ri.fields[i], cmds[i] = ri.fields[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (ri *RedisInput) PrintFields() string {
	out := ""
	for _, field := range ri.fields {
		out += field.View() + "\n"
	}
	return out
}

func (i *RedisInput) View() string {
	header := "Enter redis connection details\n"
	error := ""
	fields := i.PrintFields()

	if i.error != "" {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		error = fmt.Sprintf("\n\n%s\n%s\n",
			style.Bold(true).Render("Error"),
			style.Bold(false).Render(i.error))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, header, fields, error)
	return content
}
