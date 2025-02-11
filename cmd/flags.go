package cmd

import (
	"flag"
	"fmt"

	"github.com/mat2cc/redis_tui/tui"
)

func Run(app_version string) {
	version := flag.Bool("version", false, "Application version")
	tls := flag.Bool("tls", false, "Enable tls connection")
	tlsInsecure := flag.Bool("tls-insecure", false, "Disable tls verify")

	addressPtr := flag.String("address", "localhost:6379", "Redis server address")
	usernamePtr := flag.String("username", "", "Redis username (optional)")
	passwordPtr := flag.String("password", "", "Redis password (optional)")
	dbPtr := flag.Int("db", 0, "Redis db (default 0)")

	scanSize := flag.Int64("scan-size", 1000, "Number of keys scanned at a time")
	prettyPrintJson := flag.Bool("pp-json", true, "Pretty print JSON values")
	includeTypes := flag.Bool("include-types", true, "Include type values when querying keys. This includes a pipeline of batched TYPE commands for each key scanned, which may impact performace.")
	delimiter := flag.String("delimiter", ":", "Delimiter for key names.")

	flag.Parse()
	if *version {
		fmt.Println(app_version)
		return
	}

	tui.RunTUI(
		tui.RedisOptions{
			Address:     *addressPtr,
			Username:    *usernamePtr,
			Password:    *passwordPtr,
			DB:          *dbPtr,
			Tls:         *tls,
			TleInsecure: *tlsInsecure,
		},
		tui.ModelOptions{
			ScanSize:        *scanSize,
			PrettyPrintJson: *prettyPrintJson,
			IncludeTypes:    *includeTypes,
			Delimiter:       *delimiter,
		},
	)
}
