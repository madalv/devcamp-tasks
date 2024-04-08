package main

import (
	"adt/api"
	"adt/config"
	"adt/db"
	"adt/repository"
	"adt/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
	"net/http"
	"time"
)

func main() {
	// load cfg
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		slog.Fatal(err)
	}

	// connect to db
	conn, err := db.NewMariaDB(cfg)
	if err != nil {
		slog.Fatal(err)
	}

	// repositories
	campRepo := repository.NewCampaignRepository(conn)
	sourceRepo := repository.NewSourceRepository(conn)

	// seed db if needed
	cnt, _ := campRepo.GetCount()
	if cnt == 0 {
		seeder := util.NewSeeder(sourceRepo, campRepo)
		err = seeder.SeedDB(100, 10)
		if err != nil {
			slog.Fatal(err)
		}
	}

	// set up router
	r := chi.NewRouter()
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Logger)

	r = api.RegisterHandlers(r, campRepo, sourceRepo)

	slog.Infof("Serving HTTP on %s", cfg.HTTPort)
	err = http.ListenAndServe(cfg.HTTPort, r)
	if err != nil {
		slog.Fatal(err)
	}
}
