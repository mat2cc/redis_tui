package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Address         string
	Username        string
	Password        string
	DB              int
}

func CreateRedisClient(conn string, username string, password string, db int) (*redis.Client, error) {
	if conn == "" {
		conn = "localhost:6379"
	} else {
		conn = strings.TrimPrefix(conn, "redis://")
	}
	redis := redis.NewClient(&redis.Options{
		Addr:     conn,
		Username: username,
		Password: password,
		DB:       db,
	})
	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return redis, nil
}

func RunTUI(redis_opts RedisOptions, model_opts ModelOptions) {
	client, err := CreateRedisClient(
		redis_opts.Address,
		redis_opts.Username,
		redis_opts.Password,
		redis_opts.DB,
	)

	if err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(
		InitialModel(
			client,
            model_opts,
		),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
