package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/mat2cc/redis_tui/tui"
)

func Run(appVersion string) {
	version := flag.Bool("version", false, "Application version")

	address := flag.String("address", "localhost:6379", "Redis server address")
	username := flag.String("username", "", "Redis username (optional)")
	password := flag.String("password", "", "Redis password (optional)")
	db := flag.Int("db", 0, "Redis db (default 0)")
	tls := flag.Bool("tls", false, "Use TLS (default false)")

	scanSize := flag.Int64("scan-size", 1000, "Number of keys scanned at a time")
	prettyPrint := flag.Bool("pp-json", true, "Pretty print JSON values")
	includeTypes := flag.Bool("include-types", true, "Include type values when querying keys. This includes a pipeline of batched TYPE commands for each key scanned, which may impact performace.")
	delimiter := flag.String("delimiter", ":", "Delimiter for key names.")

	verbose := flag.Bool("verbose", false, "Verbose output")

	flag.Parse()
	if *version {
		fmt.Println(appVersion)
		return
	}

	if *verbose {
		addr, err := url.Parse(*address)
		if err != nil {
			log.Fatal(fmt.Errorf("parse address: %w", err))
		}
		log.Println("address:", addr.Redacted())
		log.Println("username:", *username)
		log.Println("password:", (*password)[:6]+"..."+(*password)[len(*password)-6:])
		log.Println("tls:", *tls)
		log.Println("db:", *db)
	}

	tui.RunTUI(
		tui.RedisOptions{
			Address:  *address,
			Username: *username,
			Password: *password,
			DB:       *db,
			TLS:      *tls,
		},
		tui.ModelOptions{
			ScanSize:        *scanSize,
			PrettyPrintJson: *prettyPrint,
			IncludeTypes:    *includeTypes,
			Delimiter:       *delimiter,
		},
	)
}
