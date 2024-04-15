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
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(h))

	// load cfg
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		slog.Error("Could not load config", "err", err)
		panic("could not load config")
	}

	// connect to db
	conn, err := db.NewMariaDB(cfg)
	if err != nil {
		slog.Error("Could not connect to db", "err", err)
		panic("could not connect to db")
	}

	// repositories
	campRepo := repository.NewCampaignRepository(conn)
	sourceRepo := repository.NewSourceRepository(conn)
	cache := repository.NewLocalCache()

	// seed db if needed
	cnt, _ := campRepo.GetCount()
	if cnt == 0 {
		seeder := util.NewSeeder(sourceRepo, campRepo)
		err = seeder.SeedDB(100, 10)
		if err != nil {
			slog.Error("Could not seed db", "err", err)
		}
	}

	// set up router
	r := chi.NewRouter()
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Logger)

	r = api.RegisterHandlers(r, campRepo, sourceRepo, cache)

	slog.Info("Serving HTTP", "port", cfg.HTTPort)
	err = http.ListenAndServe(cfg.HTTPort, r)
	if err != nil {
		slog.Error("Can't start HTTP server", "err", err)
		panic("can't start http server")
	}
}
