package usecase

import (
	"WatchHive/pkg/config"
	helper_interface "WatchHive/pkg/helper/interface"
	"WatchHive/pkg/utils/models"
	"errors"

	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
	cfg      config.Config
	helper   helper_interface.Helper
}

func NewUserUseCase(repo interfaces.UserRepository, cfg config.Config, h helper_interface.Helper) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
		cfg:      cfg,
		helper:   h,
	}
}

var InternalError = "Internal Server Error"
var ErrorHashingPassword = "Error In Hashing Password"

func (u *userUseCase) UserSignUp(user models.UserDetails) (models.TokenUsers, error) {

	userExist := u.userRepo.CheckUserAvilability(user.Email)
	if userExist {
		return models.TokenUsers{}, errors.New("user already exist, sign in")
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, errors.New("password does not match")
	}

	hashedPassword, err := u.helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New("error hashing password")
	}
	user.Password = hashedPassword
	userData, err := u.userRepo.UserSignup(user)
	if err != nil {
		return models.TokenUsers{}, err
	}
	tokenString, err := u.helper.GenerateTokenClients(userData)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token")
	}
	return models.TokenUsers{
		Users: userData,
		Token: tokenString,
	}, nil

}

func (u *userUseCase) LoginHandler(user models.UserLogin) (models.TokenUsers, error) {

	ok := u.userRepo.CheckUserAvilability(user.Email)
	if !ok {
		return models.TokenUsers{}, errors.New("the user does not exist")

	}

	isBlocked, err := u.userRepo.UserBlockStatus(user.Email)
	if err != nil {
		return models.TokenUsers{}, errors.New("internal error")
	}
	if isBlocked {
		return models.TokenUsers{}, errors.New("user is blocked")
	}
	user_details, err := u.userRepo.FindUserByEmail(user)
	if err != nil {
		return models.TokenUsers{}, errors.New("password is incorrect")
	}
	err = u.helper.CompareHashAndPassword(user_details.Password, user.Password)

	if err != nil {
		return models.TokenUsers{}, errors.New("password incorrect")
	}

	var userDetails models.UserDetailsResponse

	userDetails.Id = int(user_details.Id)
	userDetails.Name = user_details.Name
	userDetails.Email = user_details.Email
	userDetails.Phone = user_details.Phone

	tokenString, err := u.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token")
	}
	return models.TokenUsers{
		Users: userDetails,
		Token: tokenString,
	}, nil
}
