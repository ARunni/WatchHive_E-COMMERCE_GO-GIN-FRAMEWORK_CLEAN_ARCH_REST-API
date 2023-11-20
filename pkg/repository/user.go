package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {

	return &userDatabase{DB}
}

func (c *userDatabase) CheckUserAvilability(email string) bool {
	var count int
	query := fmt.Sprintf("select count (*) from users where email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (c *userDatabase) UserSignup(user models.UserDetails) (models.UserDetailsResponse, error) {

	var UserDetails models.UserDetailsResponse
	err := c.DB.Raw("insert into users (name,email,password,phone) values (?,?,?,?)returning id, name, email, phone", user.Name, user.Email, user.Password, user.Phone).Scan(&UserDetails).Error

	if err != nil {
		return models.UserDetailsResponse{}, err
	}
	return UserDetails, nil
}

func (c *userDatabase) FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error) {
	var user_details models.UserSignInResponse

	err := c.DB.Raw(`select  * from 
	users 
	where email = ?
	 and 
	 blocked = false`, user.Email).Scan(&user_details).Error

	if err != nil {
		return models.UserSignInResponse{}, errors.New("error checking user details ")
	}
	return user_details, nil

}

func (cr *userDatabase) UserBlockStatus(email string) (bool, error) {
	var isBlocked bool
	err := cr.DB.Raw("select blocked from users where email = ?", email).Scan(&isBlocked).Error
	if err != nil {
		return false, err
	}
	return isBlocked, nil
}
