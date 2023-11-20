package domain

import "WatchHive/pkg/utils/models"

type Admin struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Name     string `json:"name" gorm:"validate:required"`
	Username string `json:"email" gorm:"validate:required"`
	Password string `json:"password" gorm:"validate:required"`
}

type TockenAdmin struct {
	Admin        models.AdminDetailsResponse
	AccessToken  string
	RefreshToken string
}
