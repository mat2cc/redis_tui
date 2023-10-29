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
	input  textinput.Model
	width  int
	height int
	error  string
}

func NewRedisInput() *RedisInput {
	input := textinput.New()
	input.Placeholder = "redis://localhost:6379"
	input.Focus()
	return &RedisInput{
		input: input,
	}
}

func (i *RedisInput) Init() tea.Cmd {
	return textinput.Blink
}

func createRedisClient (conn string) (*redis.Client, error) {
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
      redis, err := createRedisClient(i.input.Value())
      if err != nil {
        i.error = err.Error()
        return i, nil
      }

			model := initialModel(redis)
			cmds := tea.Batch(setWindowSize(i.width, i.height), model.Scan())
			return model, cmds
		}
	}

	i.input, cmd = i.input.Update(msg)
	return i, cmd
}

func (i *RedisInput) View() string {
	header := "Enter redis connection string\n"
  error := ""
  if i.error != "" {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    error = fmt.Sprintf("\n\n%s\n%s\n", 
    style.Bold(true).Render("Error"),
    style.Bold(false).Render(i.error))
  }
	return header + i.input.View() + error
}
