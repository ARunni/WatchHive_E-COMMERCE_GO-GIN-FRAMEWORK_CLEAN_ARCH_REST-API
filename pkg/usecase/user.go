package usecase

import (
	"WatchHive/pkg/config"
	helper_interface "WatchHive/pkg/helper/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"

	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"
)

type userUseCase struct {
	userRepo   interfaces.UserRepository
	cfg        config.Config
	helper     helper_interface.Helper
	walletRepo interfaces.WalletRepository
}

func NewUserUseCase(repo interfaces.UserRepository, cfg config.Config, h helper_interface.Helper, wallet interfaces.WalletRepository) services.UserUseCase {
	return &userUseCase{
		userRepo:   repo,
		cfg:        cfg,
		helper:     h,
		walletRepo: wallet,
	}
}

var InternalError = "Internal Server Error"
var ErrorHashingPassword = "Error In Hashing Password"

func (u *userUseCase) UserSignUp(user models.UserDetails) (models.TokenUsers, error) {
	if user.Password == "" {
		return models.TokenUsers{}, errors.New("password cannot be empty")
	}
	if user.Name == "" {
		return models.TokenUsers{}, errors.New("name cannot be empty")
	}

	phoneNumber := u.helper.ValidatePhoneNumber(user.Phone)

	if !phoneNumber {
		return models.TokenUsers{}, errors.New("invalid phone number")
	}

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
	err = u.walletRepo.CreateWallet(userData.Id)
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
func (u *userUseCase) AddAddress(userID int, address models.AddressInfoResponse) ([]models.AddressInfoResponse, error) {

	ok, err := u.helper.ValidateAlphabets(address.Name)
	if err != nil {
		return []models.AddressInfoResponse{}, errors.New("invalid name")
	}
	if !ok {
		return []models.AddressInfoResponse{}, errors.New("invalid name")
	}
	phone := u.helper.ValidatePhoneNumber(address.Phone)
	if !phone {
		return []models.AddressInfoResponse{}, errors.New("invalid mobile number")
	}

	pin := u.helper.ValidatePin(address.Pin)
	if !pin {
		return []models.AddressInfoResponse{}, errors.New("invalid pin number")
	}

	if userID <= 0 {
		return []models.AddressInfoResponse{}, errors.New("invalid user_id")
	}

	errs := u.userRepo.CheckUserById(userID)
	if !errs {
		return []models.AddressInfoResponse{}, errors.New("user does not exist")
	}

	errResp := u.userRepo.AddAddress(userID, address)
	if errResp != nil {
		return []models.AddressInfoResponse{}, errResp
	}

	addressRep, err := u.userRepo.GetAllAddress(userID)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return addressRep, nil
}

func (u *userUseCase) ShowUserDetails(userID int) (models.UsersProfileDetails, error) {

	profile, err := u.userRepo.ShowUserDetails(userID)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return profile, nil
}

func (u *userUseCase) GetAllAddress(userID int) ([]models.AddressInfoResponse, error) {
	address, err := u.userRepo.GetAllAddress(userID)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return address, nil
}

func (u *userUseCase) EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error) {
	if user.Name == "" {
		return models.UsersProfileDetails{}, errors.New("name cannot be empty")
	}
	ok, err := u.helper.ValidateAlphabets(user.Name)
	if err != nil {
		return models.UsersProfileDetails{}, errors.New("invalid name")
	}
	if !ok {
		return models.UsersProfileDetails{}, errors.New("invalid name")
	}
	phErr := u.helper.ValidatePhoneNumber(user.Phone)
	if !phErr {
		return models.UsersProfileDetails{}, errors.New("invalid phone number")
	}
	details, err := u.userRepo.EditProfile(user)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return details, nil
}

func (u *userUseCase) ChangePassword(user models.ChangePassword) error {

	if user.NewPassWord == "" || user.ConfirmPassword == "" {
		return errors.New("password cannot be empty")
	}

	if user.NewPassWord != user.ConfirmPassword {
		return errors.New("password mismatch")
	}
	newHashed, err := u.helper.PasswordHashing(user.NewPassWord)
	if err != nil {
		return errors.New("password hashing failed")
	}

	idString := strconv.FormatUint(uint64(user.UserID), 10)

	user_details, _ := u.userRepo.FindUserById(idString)

	err = u.helper.CompareHashAndPassword(user_details.Password, user.CurrentPassWord)
	if err != nil {
		return errors.New("incorrect current password")
	}

	err = u.userRepo.ChangePassword(idString, newHashed)
	if err != nil {
		return errors.New("password cannot change")
	}
	return nil
}
