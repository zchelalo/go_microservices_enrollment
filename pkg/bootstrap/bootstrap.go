package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/zchelalo/go_microservices_domain/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Hermosillo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if os.Getenv("DB_DEBUG") == "true" {
		db = db.Debug()
	}

	if os.Getenv("DB_AUTO_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.Enrollment{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}
