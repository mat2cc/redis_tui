package tui

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Address  string
	Username string
	Password string
	DB       int
	TLS      bool
}

func CreateRedisClient(opts RedisOptions) (*redis.Client, error) {
	conn := opts.Address
	if opts.Address == "" {
		conn = "localhost:6379"
	}
	conn = strings.TrimPrefix(conn, "redis://")
	var tlsConfig *tls.Config
	if opts.TLS {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}
	redis := redis.NewClient(&redis.Options{
		Network:   "tcp",
		TLSConfig: tlsConfig,
		Addr:      conn,
		Username:  opts.Username,
		Password:  opts.Password,
		DB:        opts.DB,
	})
	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return redis, nil
}

func RunTUI(redisOpts RedisOptions, modelOpts ModelOptions) {
	client, err := CreateRedisClient(redisOpts)
	if err != nil {
		log.Fatal(fmt.Errorf("connect to redis: %v", err))
	}
	p := tea.NewProgram(
		InitialModel(client, modelOpts),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
