package api

import (
	"adt/config"
	"adt/db"
	"adt/model"
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

func TestFilterCampaigns(t *testing.T) {
	campaigns := []model.Campaign{
		{ID: 1, Blacklist: []string{}, Whitelist: []string{"cruise.com"}},
		{ID: 2, Blacklist: []string{}, Whitelist: []string{"c.com", "bruise.com"}},
		{ID: 3, Blacklist: []string{"d.com"}, Whitelist: []string{"e.com", "d.com", "example.com"}},
		{ID: 4, Blacklist: []string{"bruise.com"}, Whitelist: []string{"ise.com"}},
		{ID: 5, Blacklist: []string{}, Whitelist: []string{"sub.example.com", "y.cruise.com"}},
		{ID: 6, Blacklist: []string{"sub.test.com"}, Whitelist: []string{"sub.example.com", "m.cruise.com"}},
	}

	tests := []struct {
		name           string
		domain         string
		expectedResult []int64 // holds ids of filtered campaigns
	}{
		{"Empty domain", "", []int64{1, 2, 3, 4, 5, 6}},
		{"Domain included in whitelist", "m.CruiSe.CoM", []int64{1, 6}},
		{"Domain included in blacklist", "X.Bruise.COM", []int64{2}},
		{"Domain both blacklisted and whitelisted (blacklist takes precedence)", "d.com", []int64{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filtered := filterCampaigns(campaigns, test.domain)
			if len(filtered) != len(test.expectedResult) {
				t.Errorf("Expected %d campaigns, but got %d", len(test.expectedResult), len(filtered))
			}

			if len(filtered) > 0 {
				for i := range filtered {
					if filtered[i].ID != test.expectedResult[i] {
						t.Errorf("Expected %v campaigns, but got %v", test.expectedResult, filtered)
					}
				}
			}
		})
	}
}

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
	cache := repository.NewLocalCache()

	r := chi.NewRouter()
	r = RegisterHandlers(r, campRepo, sourceRepo, cache)

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
