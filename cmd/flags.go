package cmd

import (
	"flag"

	"github.com/mat2cc/redis_tui/tui"
)

func Run() {
	addressPtr := flag.String("address", "localhost:6379", "Redis server address")
	usernamePtr := flag.String("username", "", "Redis username (optional)")
	passwordPtr := flag.String("password", "", "Redis password (optional)")
	dbPtr := flag.Int("db", 0, "Redis db, defaults to 0 (optional)")
    scanSize := flag.Int64("scan-size", 1000, "Number of keys scanned at a time")

	flag.Parse()

	tui.RunTUI(*addressPtr, *usernamePtr, *passwordPtr, *dbPtr, *scanSize)
}
