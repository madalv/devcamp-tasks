package model

type Source struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type CreateSourceDTO struct {
	Name        string  `json:"name"`
	CampaignIDs []int64 `json:"campaign_ids"`
}
