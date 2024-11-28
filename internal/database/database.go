package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func Connect(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDatabase(_db *gorm.DB) *Database {
	return &Database{
		db: _db,
	}
}
