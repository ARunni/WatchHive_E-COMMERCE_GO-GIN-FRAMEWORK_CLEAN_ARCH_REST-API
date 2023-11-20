package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"

	"gorm.io/gorm"
)
type adminRepository struct {
	DB *gorm.DB
}
func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository  {
	return &adminRepository{
		DB: DB,
	}
}

func (ad *adminRepository) LoginHandler(adminDetails models.AdminLogin)(domain.Admin,error)  {

	var adminCompareDetails domain.Admin
	if err := ad.DB.Raw("se;ect * from admins where username = ?",adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.Admin{},err
	}
	return adminCompareDetails,nil
}