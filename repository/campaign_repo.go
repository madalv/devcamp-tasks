package repository

import (
	"adt/model"
	"github.com/jmoiron/sqlx"
)

type CampaignRepository struct {
	db *sqlx.DB
}

func NewCampaignRepository(db *sqlx.DB) *CampaignRepository {
	return &CampaignRepository{db: db}
}

func (r *CampaignRepository) Create(c *model.CreateCampaignDTO) (campID int64, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	res, err := tx.Exec(`insert into campaigns (name) values (?)`, c.Name)
	if err != nil {
		return
	}
	campID, err = res.LastInsertId()
	if err != nil {
		return
	}

	if len(c.SourceIDs) > 0 {
		for _, sourceID := range c.SourceIDs {
			_, err = tx.Exec(`insert into campaigns_sources (source_id, campaign_id) values (?, ?)`, sourceID, campID)
			if err != nil {
				return
			}
		}
	}
	return
}

func (r *CampaignRepository) GetAllNoSources() (camps []model.Campaign, err error) {
	query :=
		`select c.name, c.id from campaigns c
		left join campaigns_sources cs on c.id = cs.campaign_id
		where cs.source_id is null`
	err = r.db.Select(&camps, query)
	return
}

func (r *CampaignRepository) Update(c *model.Campaign) (err error) {
	_, err = r.db.Exec(`update campaigns set name = ? where id = ?`, c.Name, c.ID)
	return
}

func (r *CampaignRepository) Delete(id int64) (err error) {
	_, err = r.db.Exec(`update campaigns set name = ? where id = ?`, id)
	return
}
