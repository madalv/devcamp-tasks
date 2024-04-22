package model

import (
	"github.com/brianvoe/gofakeit/v7"
	"time"
)

const (
	WHITELIST = "white"
	BLACKLIST = "black"
)

type Campaign struct {
	ID         int64               `db:"id" json:"id"`
	Name       string              `db:"name" json:"name"`
	DomainList map[string]struct{} `db:"blacklist" json:"domain_list"`
	ListType   string              `db:"list_type" json:"list_type"`
}

func (c *Campaign) Call() int {
	time.Sleep(1 * time.Second)
	return gofakeit.Number(0, 100)
}

type CreateCampaignDTO struct {
	Name       string   `json:"name"`
	SourceIDs  []int64  `json:"source_ids"`
	DomainList []string `json:"domain_list"`
	ListType   string   `json:"list_type"`
}
