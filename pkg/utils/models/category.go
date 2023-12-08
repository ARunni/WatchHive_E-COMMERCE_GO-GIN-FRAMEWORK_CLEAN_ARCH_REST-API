package models

type CategoryResp struct {
	Id       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category" gorm:"unique;not null"`
}

type SetNewName struct {
	CurrentId int    `json:"current_id"`
	New       string `json:"new"`
}
