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
	conn textinput.Model
	username   textinput.Model
	password   textinput.Model
	width      int
	height     int
	error      string
}

func NewRedisInput() *RedisInput {
	connection := textinput.New()
	connection.Placeholder = "redis://localhost:6379"
	connection.Focus()

	return &RedisInput{
		conn: connection,
		username: textinput.New(),
		password: textinput.New(),
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
			redis, err := createRedisClient(i.conn.Value())
			if err != nil {
				i.error = err.Error()
				return i, nil
			}

			model := initialModel(redis)
			cmds := tea.Batch(setWindowSize(i.width, i.height), model.Scan())
			return model, cmds
		}
	}

	i.conn, cmd = i.conn.Update(msg)
	return i, cmd
}

func (i *RedisInput) View() string {
	header := "Enter redis connection string\n"
    fields := fmt.Sprintf("Username: %s\nPassword: %s\n", i.username.View(), i.password.View())
	error := ""
	if i.error != "" {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
		error = fmt.Sprintf("\n\n%s\n%s\n",
			style.Bold(true).Render("Error"),
			style.Bold(false).Render(i.error))
	}

	return header + i.conn.View() + error
}
