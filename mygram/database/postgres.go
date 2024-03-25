package database

import (
	"fmt"
	"github.com/refandas/scalable-web-service/mygram/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Postgres struct {
	DB  *gorm.DB
	Err error
}

func NewPostgres() *Postgres {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), nil)

	err = db.Debug().AutoMigrate(&core.Comment{}, &core.Photo{}, &core.SocialMedia{}, &core.User{})
	if err != nil {
		return nil
	}

	return &Postgres{
		DB:  db,
		Err: err,
	}
}
