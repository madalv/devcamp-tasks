package model

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

type CreateCampaignDTO struct {
	Name       string   `json:"name"`
	SourceIDs  []int64  `json:"source_ids"`
	DomainList []string `json:"domain_list"`
	ListType   string   `json:"list_type"`
}
