package api

import (
	"adt/model"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type CampaignRepo interface {
	Create(c *model.CreateCampaignDTO) (campID int64, err error)
	GetAllBySourceID(sourceID int) (camps []model.Campaign, err error)
}

type SourceRepo interface {
	Create(c *model.CreateSourceDTO) (sourceID int64, err error)
}

type Cache interface {
	Get(key string) ([]model.Campaign, bool)
	Put(key string, value []model.Campaign, ttl time.Duration)
}

type handler struct {
	campaignRepo CampaignRepo
	sourceRepo   SourceRepo
	cache        Cache
}

func RegisterHandlers(r *chi.Mux, cr CampaignRepo, sr SourceRepo, cache Cache) *chi.Mux {
	h := &handler{
		campaignRepo: cr,
		sourceRepo:   sr,
		cache:        cache}

	r.Route("/api/sources", func(r chi.Router) {
		r.Get("/{sourceID}/campaigns", h.getCampaignsForSource)
	})

	return r
}

func (h *handler) getCampaignsForSource(w http.ResponseWriter, r *http.Request) {
	sourceIDStr := chi.URLParam(r, "sourceID")
	sourceID, err := strconv.Atoi(sourceIDStr)
	if err != nil {
		http.Error(w, "invalid sourceID", http.StatusBadRequest)
		return
	}
	domain := strings.ToLower(r.URL.Query().Get("domain"))
	var filteredCamps []model.Campaign

	// try to retrieve from cache first
	camps, ok := h.cache.Get(fmt.Sprintf("CAMPS_FOR_SRC_%d", sourceID))
	if ok {
		filteredCamps = filterCampaigns(camps, domain)

		encoded, err := json.Marshal(filteredCamps)
		if err != nil {
			http.Error(w, "could not marshal response", http.StatusInternalServerError)
			return
		}

		writeJson(w, encoded)
		return
	}

	camps, err = h.campaignRepo.GetAllBySourceID(sourceID)
	if err != nil {
		slog.Error(err)
	}
	filteredCamps = filterCampaigns(camps, domain)

	encoded, err := json.Marshal(filteredCamps)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	// set to cache
	h.cache.Put(fmt.Sprintf("CAMPS_FOR_SRC_%d", sourceID), camps, time.Second*5)
	writeJson(w, encoded)
}

func filterCampaigns(camps []model.Campaign, domain string) (filtered []model.Campaign) {
	if domain == "" {
		return camps
	}
	for _, c := range camps {
		if !slices.Contains(c.Blacklist, domain) &&
			slices.Contains(c.Whitelist, domain) {
			filtered = append(filtered, c)
		}
	}
	return
}

func writeJson(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
