package repository

import (
	"adt/model"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"strings"
)

type CampaignRepository struct {
	db *sqlx.DB
}

func NewCampaignRepository(db *sqlx.DB) *CampaignRepository {
	slog.Info("Setting up new Campaign Repository . . .")
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

	res, err := tx.Exec(`insert into campaigns (name, domain_list, list_type) values (?, ?, ?)`,
		c.Name,
		strings.ToLower(strings.Join(c.DomainList, ",")),
		c.ListType,
	)

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

func (r *CampaignRepository) GetAllBySourceID(sourceID int) (camps []model.Campaign, err error) {
	query :=
		`select c.name, c.id, c.domain_list, c.list_type from campaigns c
		join campaigns_sources cs on c.id = cs.campaign_id
		where cs.source_id = ?`
	rows, err := r.db.Queryx(query, sourceID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var camp model.Campaign
		var listStr string
		err = rows.Scan(&camp.Name, &camp.ID, &listStr, &camp.ListType)
		if err != nil {
			return
		}

		camp.DomainList = strings.Split(listStr, ",")
		camps = append(camps, camp)
	}

	return
}

func (r *CampaignRepository) GetCount() (count int, err error) {
	err = r.db.Get(&count, `select count(*) from campaigns`)
	return
}
