package main

import (
	"adt/config"
	"adt/db"
	"adt/repository"
	"adt/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
)

func main() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		slog.Fatal(err)
	}

	conn, err := db.NewMariaDB(cfg)
	if err != nil {
		slog.Fatal(err)
	}

	campRepo := repository.NewCampaignRepository(conn)
	sourceRepo := repository.NewSourceRepository(conn)

	seeder := util.NewSeeder(sourceRepo, campRepo)
	err = seeder.SeedDB(100, 10)
	if err != nil {
		slog.Fatal(err)
	}
}
