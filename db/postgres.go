package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct{}

func (p *Postgres) ConnectUrl(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *Postgres) Reset(db *gorm.DB, table string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE " + table).Error; err != nil {
			return err
		}

		if err := tx.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1").Error; err != nil {
			return err
		}

		return nil
	})
}
