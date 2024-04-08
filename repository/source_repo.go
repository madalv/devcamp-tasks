package repository

import (
	"adt/model"
	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
)

type SourceRepository struct {
	db *sqlx.DB
}

func NewSourceRepository(db *sqlx.DB) *SourceRepository {
	slog.Info("Setting up new Source Repository . . .")
	return &SourceRepository{db: db}
}

func (r *SourceRepository) Create(c *model.CreateSourceDTO) (sourceID int64, err error) {
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

	res, err := tx.Exec(`insert into sources (name) values (?)`, c.Name)
	if err != nil {
		return
	}
	sourceID, err = res.LastInsertId()
	if err != nil {
		return
	}

	if len(c.CampaignIDs) > 0 {
		for _, campID := range c.CampaignIDs {
			_, err = tx.Exec(`insert into campaigns_sources(source_id, campaign_id) values (?, ?)`, sourceID, campID)
			if err != nil {
				return
			}
		}
	}
	return
}

func (r *SourceRepository) GetSourcesByCampNr(limit uint) (camps []model.Source, err error) {
	query :=
		`select s.name, s.id, count(cs.campaign_id) from sources s
		left join campaigns_sources cs on s.id = cs.source_id
		group by s.id
		order by count(cs.campaign_id) desc
		limit ?;`
	err = r.db.Select(&camps, query, limit)
	return
}

func (r *SourceRepository) Update(c *model.Source) (err error) {
	_, err = r.db.Exec(`update sources set name = ? where id = ?`, c.Name, c.ID)
	return
}

func (r *SourceRepository) Delete(id int64) (err error) {
	_, err = r.db.Exec(`update sources set name = ? where id = ?`, id)
	return
}
