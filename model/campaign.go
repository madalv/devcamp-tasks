package model

type Campaign struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type CreateCampaignDTO struct {
	Name      string  `json:"name"`
	SourceIDs []int64 `json:"source_ids"`
}
