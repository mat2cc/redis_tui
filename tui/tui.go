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
	Address     string
	Username    string
	Password    string
	DB          int
	Tls         bool
	TleInsecure bool
}

func CreateRedisClient(conn string, username string, password string, db int, enableTls, tlsInsecure bool) (*redis.Client, error) {
	if conn == "" {
		conn = "localhost:6379"
	}
	if !strings.HasPrefix(conn, "redis://") && !strings.HasPrefix(conn, "rediss://") {
		if enableTls {
			conn = "rediss://" + conn
		} else {
			conn = "redis://" + conn
		}
	}
	options, err := redis.ParseURL(conn)
	if err != nil {
		return nil, err
	}
	options.Username = username
	options.Password = password
	if enableTls && tlsInsecure {
		options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	redis := redis.NewClient(options)
	_, err = redis.Ping(context.Background()).Result()
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
		redis_opts.Tls,
		redis_opts.TleInsecure,
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
