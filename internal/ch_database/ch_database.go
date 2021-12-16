package ch_database

import (
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/history/internal/config"
	"github.com/rusrafkasimov/history/pkg/models"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

func InitializeDB(config *config.Configuration, logger promtail.Client) (*gorm.DB, error) {
	dbDsn, err := config.Get("CLICKHOUSE_DSN")
	if err != nil {
		logger.Errorf("Error: can't parse db dsn. %s", err.Error())
		return nil, err
	}

	db, err := gorm.Open(clickhouse.Open(dbDsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("Error: can't click house. %s", err.Error())
		return nil, err
	}

	logger.Infof("We are click house :)")

	return db, err
}

func MigrateHistoryModels (db *gorm.DB) {
	err := db.Debug().AutoMigrate(
		&models.AccountHistory{},
		)
	if err != nil {
		fmt.Println("can't migrate models to DB")
		return
	}
}
