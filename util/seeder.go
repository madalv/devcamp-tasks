package util

import (
	"adt/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gookit/slog"
)

type SourceRepo interface {
	Create(c *model.CreateSourceDTO) (sourceID int64, err error)
}

type CampaignRepo interface {
	Create(c *model.CreateCampaignDTO) (campID int64, err error)
}

type DBSeeder struct {
	sourceRepo   SourceRepo
	campaignRepo CampaignRepo
}

func NewSeeder(sr SourceRepo, cr CampaignRepo) *DBSeeder {
	return &DBSeeder{
		campaignRepo: cr,
		sourceRepo:   sr}
}

func (s *DBSeeder) SeedDB(rows, maxCampPerSource int) error {
	slog.Info("Seeding DB . . . ")
	campIDs, err := s.seedCampaigns(rows)
	if err != nil {
		return err
	}
	slog.Info("Seeded campaigns table")

	err = s.seedSources(campIDs, rows, maxCampPerSource)
	if err != nil {
		return err
	}
	slog.Info("Seeded sources table")
	return nil
}

func (s *DBSeeder) seedCampaigns(rows int) ([]int64, error) {
	ids := make([]int64, rows)
	for i := 0; i < rows; i++ {
		id, err := s.campaignRepo.Create(&model.CreateCampaignDTO{
			Name:      gofakeit.Word(),
			SourceIDs: nil,
		})
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (s *DBSeeder) seedSources(campIDs []int64, rows, maxCampPerSource int) error {
	for i := 0; i < rows; i++ {
		nrCampPerSource := gofakeit.Number(0, maxCampPerSource)
		source := model.CreateSourceDTO{
			Name:        gofakeit.DomainName(),
			CampaignIDs: make([]int64, nrCampPerSource),
		}

		for j := 0; j < nrCampPerSource; {
			campID := campIDs[gofakeit.IntN(rows)]
			if !contains(source.CampaignIDs, campID) {
				source.CampaignIDs[j] = campID
				j++
			}
		}
		_, err := s.sourceRepo.Create(&source)
		if err != nil {
			return err
		}
	}
	return nil
}

func contains(slice []int64, value int64) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}