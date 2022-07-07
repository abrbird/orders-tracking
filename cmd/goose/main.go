package main

import (
	"flag"
	"github.com/pressly/goose/v3"
	"log"
	"os"

	"github.com/abrbird/orders-tracking/config"
	_ "github.com/abrbird/orders-tracking/migrations"
	_ "github.com/lib/pq"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("goose: failed on args parsing")
	}

	args := flags.Args()
	if len(args) < 1 {
		flags.Usage()
		log.Fatal("goose: not enough args")
	}

	dbString, command := cfg.Database.String(), args[1]

	db, err := goose.OpenDBWithDriver(cfg.Database.DBMS, dbString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
