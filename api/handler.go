package api

import (
	"adt/model"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CampaignRepo interface {
	Create(c *model.CreateCampaignDTO) (campID int64, err error)
	GetAllBySourceID(sourceID int) (camps []model.Campaign, err error)
}

type SourceRepo interface {
	Create(c *model.CreateSourceDTO) (sourceID int64, err error)
}

type handler struct {
	campaignRepo CampaignRepo
	sourceRepo   SourceRepo
}

func RegisterHandlers(r *chi.Mux, cr CampaignRepo, sr SourceRepo) *chi.Mux {
	h := &handler{
		campaignRepo: cr,
		sourceRepo:   sr}
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

	camps, err := h.campaignRepo.GetAllBySourceID(sourceID)

	jsonData, err := json.Marshal(camps)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
