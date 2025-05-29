package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormPostgres interface {
	GetConnection() *gorm.DB
}

type gormPostgresImpl struct {
	master *gorm.DB
}

func NewGormPostgres() GormPostgres {
	return &gormPostgresImpl{
		master: connect(),
	}
}

func connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL_CONFIGURATION")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func (g *gormPostgresImpl) GetConnection() *gorm.DB {
	return g.master
}
