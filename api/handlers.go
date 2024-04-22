package api

import (
	"adt/model"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	domain := r.URL.Query().Get("domain")
	var filteredCamps []model.Campaign

	// try to retrieve from cache first
	camps, ok := h.cache.Get(fmt.Sprintf("CAMPS_FOR_SRC_%d", sourceID))
	if ok {
		slog.Info("got from cache", "camps", camps)
		filteredCamps = filterCampaigns(camps, domain)
		sortCampaignsByBid(filteredCamps)

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
		http.Error(w, "could not get campaigns", http.StatusInternalServerError)
		return
	}
	filteredCamps = filterCampaigns(camps, domain)
	sortCampaignsByBid(filteredCamps)

	encoded, err := json.Marshal(filteredCamps)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}
	// set to cache
	h.cache.Put(fmt.Sprintf("CAMPS_FOR_SRC_%d", sourceID), camps, time.Second*10)
	writeJson(w, encoded)
}

type bid struct {
	index  int
	amount int
}

func sortCampaignsByBid(camps []model.Campaign) {
	n := len(camps)
	var wg sync.WaitGroup
	bidChan := make(chan *bid)
	bids := make([]int, n)
	wg.Add(n)

	for i := range camps {
		go func(i int) {
			c := &camps[i]
			amount := c.Call()
			bidChan <- &bid{i, amount}
		}(i)
	}

	go func() {
		for bid := range bidChan {
			bids[bid.index] = bid.amount
			wg.Done()
		}
	}()

	wg.Wait()
	close(bidChan)

	// sort slice according to bids
	sort.Slice(camps, func(i, j int) bool {
		return bids[i] > bids[j]
	})
}

func filterCampaigns(camps []model.Campaign, domain string) (filtered []model.Campaign) {
	if domain == "" {
		return camps
	} else {
		domain = strings.ToLower(domain)
	}

	for _, c := range camps {
		contained := domainInList(domain, c.DomainList)

		if (contained && c.ListType == model.BLACKLIST) ||
			(!contained && c.ListType == model.WHITELIST) {
			//slog.Debug("Campaign skipped", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
			continue
		}
		//slog.Debug("Campaign good", "cid", c.ID, "contained?", contained, "domain", domain, "type", c.ListType, "list", c.DomainList)
		filtered = append(filtered, c)
	}
	return
}

/*
I make the assumption that if a camp. has domain "a.com" in its whitelist/blacklist,
the queryDomain "b.a.com" will include/filter out the camp., but the opposite is not true, i.e.
if a camp. has a domain "c.a.com" it its whitelist/blacklist, the queryDomain "a.com" will not
either include/filter out the campaign.
*/
func domainInList(queryDomain string, dMap map[string]struct{}) bool {
	parts := strings.Split(queryDomain, ".")
	for i := 0; i < len(parts)-1; i++ {
		currDomain := strings.Join(parts[i:], ".")
		_, ok := dMap[currDomain]
		if ok {
			return true
		}
	}

	return false
}

func writeJson(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
