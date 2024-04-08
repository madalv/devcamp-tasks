package api

import (
	"adt/config"
	"adt/db"
	"adt/repository"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGetCampaignsForSource(b *testing.B) {
	cfg, err := config.LoadConfig("../config.yaml")
	if err != nil {
		slog.Fatal(err)
	}

	conn, err := db.NewMariaDB(cfg)
	if err != nil {
		slog.Fatal(err)
	}

	campRepo := repository.NewCampaignRepository(conn)
	sourceRepo := repository.NewSourceRepository(conn)

	r := chi.NewRouter()
	r = RegisterHandlers(r, campRepo, sourceRepo)

	rr := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("/api/sources/%d/campaigns", gofakeit.Number(1, 100))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			b.Fatal(err)
		}
		r.ServeHTTP(rr, req)
	}
}
