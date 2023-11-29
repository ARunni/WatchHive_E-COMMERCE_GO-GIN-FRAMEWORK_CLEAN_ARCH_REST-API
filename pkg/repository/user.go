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

	result := c.DB.Raw("INSERT INTO users (name, email, password, phone) VALUES (?, ?, ?, ?) RETURNING id, name, email, phone",
		user.Name, user.Email, user.Password, user.Phone).Scan(&UserDetails)

	if result.Error != nil {
		return UserDetails, result.Error
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

func (c *userDatabase) FindUserById(id string) (models.UserSignInResponse, error) {
	var user_details models.UserSignInResponse
	err := c.DB.Raw("select * from users where id = ?", id).Scan(&user_details).Error
	if err != nil {
		return models.UserSignInResponse{}, errors.New("error in checking user details")
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

func (db *userDatabase) CheckIfUserAddress(userID int) bool {
	var count int
	qurry := "select count(*) from addresses where user_id = $1"
	if err := db.DB.Raw(qurry, userID).Scan(&count).Error; err != nil {
		return false
	}
	return true

}

func (db *userDatabase) CheckUserById(userID int) bool {
	var count int
	qurry := "select count(*) from users where id = $1"
	if err := db.DB.Raw(qurry, userID).Scan(&count).Error; err != nil {
		return false
	}
	return true
}

func (db *userDatabase) AddAddress(userID int, address models.AddressInfoResponse) (models.AddressInfoResponse, error) {

	querry := "INSERT INTO addresses(user_id,name,house_name,street,city,state,phone,pin) VALUES (?,?,?,?,?,?,?,?)"
	err := db.DB.Exec(querry, userID, address.Name, address.HouseName, address.Street, address.City, address.State, address.Phone, address.Pin).Error
	if err != nil {
		return models.AddressInfoResponse{}, errors.New("could not add address, db error")
	}
	return models.AddressInfoResponse{}, nil

}

func (db *userDatabase) ShowUserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	query := "SELECT id,name,email,phone from users where id = ?"

	result := db.DB.Raw(query, userID).Scan(&userDetails)
	if result.Error != nil {
		return models.UsersProfileDetails{}, result.Error
	}
	return userDetails, nil

}

func (db *userDatabase) GetAllAddress(userID int) ([]models.AddressInfoResponse, error) {
	qurey := "SELECT * from addresses where user_id = ?"
	var address []models.AddressInfoResponse
	result := db.DB.Raw(qurey, userID).Scan(&address)
	if result.Error != nil {
		return []models.AddressInfoResponse{}, result.Error
	}
	return address, nil
}

func (db *userDatabase) EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error) {
	querry := "UPDATE  users SET name = ?, email = ?, phone = ? WHERE id = ?"
	err := db.DB.Exec(querry, user.Name, user.Email, user.Phone, user.ID).Error
	if err != nil {
		return models.UsersProfileDetails{}, err
	}

	return user, nil
}

func (db *userDatabase) ChangePassword(userID, password string) error {
	query := "update users set password = ? where id = ?"
	err := db.DB.Exec(query, password, userID).Error
	if err != nil {
		return err
	}
	return nil
}
