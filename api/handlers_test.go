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
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func BenchFilterCampaigns() {

}

func TestFilterCampaigns(t *testing.T) {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(h))

	campaigns := []model.Campaign{
		{ID: 1, DomainList: []string{"123domain.xyz", "example.com"}, ListType: model.WHITELIST},
		{ID: 2, DomainList: []string{"123test.abc", "random.com"}, ListType: model.BLACKLIST},
		{ID: 3, DomainList: []string{"main.xyz", "123domain.xy"}, ListType: model.WHITELIST},
		{ID: 4, DomainList: []string{"example.com", "est.abc", "test.ab"}, ListType: model.BLACKLIST},
	}

	tests := []struct {
		name           string
		domain         string
		expectedResult []int64 // holds ids of filtered campaigns
	}{
		{"Empty domain", "", []int64{1, 2, 3, 4}},
		{"Domain included in whitelist", "exaMple.com", []int64{1, 2}},
		{"Domain included in blacklist", "RANDom.cOm", []int64{4}},
		{"Subdomain included in whitelist", "x.123domain.xyz", []int64{1, 2, 4}},
		{"Subdomain included in blacklist", "x.123test.abc", []int64{4}},
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
		slog.Error("err", err)
		return
	}

	conn, err := db.NewMariaDB(cfg)
	if err != nil {
		slog.Error("err", err)
		return
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
