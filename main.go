package main

import (
	"adt/config"
	db "adt/db/sqlc"
	"database/sql"

	"github.com/gookit/slog"
)

func main() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		slog.Fatal(err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		slog.Fatal(err)
	}

	querier := db.New(conn)

	seeder := db.NewSeeder(querier)

	err = seeder.SeedDB(100)
	if err != nil {
		slog.Error(err)
	}
}
