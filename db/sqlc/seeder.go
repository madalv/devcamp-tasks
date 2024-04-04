package db

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
)

type Seeder struct {
	querier Querier
	ctx     context.Context
}

func NewSeeder(q Querier) *Seeder {
	return &Seeder{
		querier: q,
		ctx:     context.Background(),
	}
}

func (s *Seeder) SeedDB(nrRows int) error {
	err := s.generateCampaigns(nrRows)
	if err != nil {
		return err
	}
	err = s.generateSources(nrRows)
	if err != nil {
		return err
	}

	return nil
}

func (s *Seeder) generateCampaigns(n int) error {
	for i := 0; i < n; i++ {
		err := s.querier.CreateCampaign(s.ctx, gofakeit.Word())
		if err != nil {
			return err
		}
	}
	slog.Info("Campaigns generated")
	return nil
}

func (s *Seeder) generateSources(n int) error {
	for i := 0; i < n; i++ {
		err := s.querier.CreateSource(s.ctx, gofakeit.DomainName())
		if err != nil {
			return err
		}

		err = s.generateSourceCampaignLink(n, i+1, 10)
		if err != nil {
			return err
		}
	}
	slog.Info("Sources generated")
	return nil
}

func (s *Seeder) generateSourceCampaignLink(n, sourceId, campPerSource int) error {
	nr := gofakeit.Number(0, campPerSource)
	for i := 0; i < nr; i++ {
		err := s.querier.CreateCampaignSourceLink(s.ctx, &CreateCampaignSourceLinkParams{
			SourceID:   uint64(sourceId),
			CampaignID: uint64(gofakeit.Number(1, n)),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
