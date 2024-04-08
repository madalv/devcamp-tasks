package api

import (
	"adt/model"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
	Get(key string) ([]byte, bool)
	Put(key string, value []byte, ttl time.Duration)
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

	// try to retrieve from cache first
	data, ok := h.cache.Get(fmt.Sprintf("CAMPS_FOR_SRC_%d", sourceID))
	if ok {
		writeJson(w, data)
		return
	}

	camps, err := h.campaignRepo.GetAllBySourceID(sourceID)

	data, err = json.Marshal(camps)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	// set to cache
	h.cache.Put(fmt.Sprintf("CAMPS_FOR_SRC_%d", sourceID), data, time.Second*5)

	writeJson(w, data)
}

func writeJson(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
