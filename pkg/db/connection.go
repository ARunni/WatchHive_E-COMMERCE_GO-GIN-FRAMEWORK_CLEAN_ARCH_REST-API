package db

import (
	"WatchHive/pkg/config"
	"WatchHive/pkg/domain"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{SkipDefaultTransaction: true})

	if err := db.AutoMigrate(&domain.Users{}); err != nil {
		return db, err
	}

	if err := db.AutoMigrate(&domain.Admin{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Category{}); err != nil {
		return db, err
	}
	if err := db.AutoMigrate(&domain.Inventory{}); err != nil {
		return db, err
	}

	CheckAndCreateAdmin(db)

	return db, dbErr
}

func CheckAndCreateAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.Admin{}).Count(&count)
	if count == 0 {
		password := "1234"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		admin := domain.Admin{
			ID:       1,
			Name:     "watchHive",
			Username: "watchhive@watchhive.com",
			Password: string(hashedPassword),
		}
		db.Create(&admin)
	}
}
