package cmd

import (
	"flag"

	"github.com/mat2cc/redis_tui/tui"
)

type Options struct {
	address         string
	username        string
	password        string
	db              int
	scanSize        int64
	prettyPrintJson bool
	includeTypes    bool
}

func Run() {
	addressPtr := flag.String("address", "localhost:6379", "Redis server address")

	usernamePtr := flag.String("username", "", "Redis username (optional)")
	passwordPtr := flag.String("password", "", "Redis password (optional)")
	dbPtr := flag.Int("db", 0, "Redis db (default 0)")

	scanSize := flag.Int64("scan-size", 1000, "Number of keys scanned at a time")
	prettyPrintJson := flag.Bool("pp-json", true, "Pretty print JSON values")
	includeTypes := flag.Bool("include-types", true, "Include type values when querying keys. This includes a pipeline of batched TYPE commands for each key scanned, which may impact performace.")

	flag.Parse()

	tui.RunTUI(tui.Options{
		Address:         *addressPtr,
		Username:        *usernamePtr,
		Password:        *passwordPtr,
		DB:              *dbPtr,
		ScanSize:        *scanSize,
		PrettyPrintJson: *prettyPrintJson,
		IncludeTypes:    *includeTypes,
	})
}
