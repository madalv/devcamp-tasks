package model

type Campaign struct {
	ID        int64    `db:"id" json:"id"`
	Name      string   `db:"name" json:"name"`
	Blacklist []string `db:"blacklist" json:"blacklist"`
	Whitelist []string `db:"whitelist" json:"whitelist"`
}

type CreateCampaignDTO struct {
	Name      string   `json:"name"`
	SourceIDs []int64  `json:"source_ids"`
	Blacklist []string `json:"blacklist"`
	Whitelist []string `json:"whitelist"`
}
