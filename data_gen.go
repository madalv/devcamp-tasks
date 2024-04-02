package main

import (
	"database/sql"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/db")
	if err != nil {
		slog.Fatal(err)
	}
	defer db.Close()

	n := 100

	err = generateCampaigns(db, n)
	if err != nil {
		slog.Fatal(err)
	}
	err = generateSources(db, n)
	if err != nil {
		slog.Fatal(err)
	}
}

func generateCampaigns(db *sql.DB, n int) error {
	for i := 0; i < n; i++ {
		_, err := db.Exec("INSERT INTO campaigns (name) VALUES (?)", gofakeit.Noun())
		if err != nil {
			return err
		}
	}
	slog.Info("Campaigns generated")
	return nil
}

func generateSources(db *sql.DB, n int) error {
	for i := 0; i < n; i++ {
		res, err := db.Exec("INSERT INTO sources (name) VALUES (?) ", gofakeit.DomainName())
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		err = generateSourceCampaignLink(db, n, int(id), 10)
		if err != nil {
			return err
		}
	}
	slog.Info("Sources generated")
	return nil
}

func generateSourceCampaignLink(db *sql.DB, n, source_id, campPerSource int) error {
	nr := gofakeit.Number(0, campPerSource)
	for i := 0; i < nr; i++ {
		_, err := db.Exec("INSERT INTO campaigns_sources (source_id, campaign_id) VALUES (?, ?)", source_id, gofakeit.Number(1, n))
		if err != nil {
			return err
		}
	}
	return nil
}
