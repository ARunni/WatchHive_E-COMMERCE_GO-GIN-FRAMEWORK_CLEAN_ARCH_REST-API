package usecase

import (
	"WatchHive/pkg/config"
	helper_interface "WatchHive/pkg/helper/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"

	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"

	"github.com/google/uuid"
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
		return models.TokenUsers{}, errors.New("password " + errmsg.ErrFieldEmpty)
	}
	if user.Name == "" {
		return models.TokenUsers{}, errors.New("name " + errmsg.ErrFieldEmpty)
	}

	phoneNumber := u.helper.ValidatePhoneNumber(user.Phone)

	if !phoneNumber {
		return models.TokenUsers{}, errors.New(errmsg.ErrInvalidPhone)
	}

	userExist := u.userRepo.CheckUserAvilability(user.Email)
	if userExist {
		return models.TokenUsers{}, errors.New(errmsg.ErrAlreadyUser)
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, errors.New(errmsg.ErrPasswordMatch)
	}

	hashedPassword, err := u.helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrPasswordHash)
	}
	user.Password = hashedPassword
	userData, err := u.userRepo.UserSignup(user)
	if err != nil {
		return models.TokenUsers{}, err
	}
	id := uuid.New().ID()
	str := strconv.Itoa(int(id))
	userReferral := str[:8]
	err = u.userRepo.NewReferralEntry(userData.Id, userReferral)
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrCreateRefferal)
	}
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrWriteDB)
	}

	err = u.walletRepo.CreateWallet(userData.Id)
	if err != nil {
		return models.TokenUsers{}, err
	}

	//
	if user.ReferralCode != "" {
		// first check whether if a user with that referralCode exist
		referredId, err := u.userRepo.GetUserIdFromReferralCode(user.ReferralCode)
		if err != nil {
			return models.TokenUsers{}, errors.New(errmsg.ErrGetDB)
		}
		if referredId != 0 {
			referralAmount := 150
			err := u.userRepo.UpdateReferralAmount(float64(referralAmount), referredId, userData.Id)
			if err != nil {
				return models.TokenUsers{}, err
			}
			// referreason := "Amount credited for used referral code"
			// err = u.userRepo.UpdateHistory(userData.Id, 0, float64(referralAmount), referreason)
			// if err != nil {
			// 	return models.TokenUsers{}, err
			// }
			amount, err := u.userRepo.AmountInRefferals(userData.Id)
			if err != nil {
				return models.TokenUsers{}, err
			}
			wallectExist, err := u.walletRepo.IsWalletExist(referredId)
			if err != nil {
				return models.TokenUsers{}, err
			}
			if !wallectExist {
				err = u.walletRepo.CreateWallet(referredId)
				if err != nil {
					return models.TokenUsers{}, err
				}
			}
			err = u.walletRepo.AddToWallet(referredId, amount)
			if err != nil {
				return models.TokenUsers{}, err
			}
			err = u.walletRepo.AddToWallet(userData.Id, float64(referralAmount))
			if err != nil {
				return models.TokenUsers{}, err
			}
			// reason := "Amount credited for refer a new person"
			// err = u.userRepo.UpdateHistory(referredId, 0, amount, reason)
			// if err != nil {
			// 	return models.TokenUsers{}, err
			// }
		}
	}
	//
	tokenString, err := u.helper.GenerateTokenClients(userData)
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrCreateTocken)
	}
	return models.TokenUsers{
		Users: userData,
		Token: tokenString,
	}, nil

}

func (u *userUseCase) LoginHandler(user models.UserLogin) (models.TokenUsers, error) {

	ok := u.userRepo.CheckUserAvilability(user.Email)
	if !ok {
		return models.TokenUsers{}, errors.New(errmsg.ErrUserExistFalse)

	}

	isBlocked, err := u.userRepo.UserBlockStatus(user.Email)
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrInternal)
	}
	if isBlocked {
		return models.TokenUsers{}, errors.New(errmsg.ErrUserBlockTrue)
	}
	user_details, err := u.userRepo.FindUserByEmail(user)
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrPassword)
	}
	err = u.helper.CompareHashAndPassword(user_details.Password, user.Password)

	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrPassword)
	}

	var userDetails models.UserDetailsResponse

	userDetails.Id = int(user_details.Id)
	userDetails.Name = user_details.Name
	userDetails.Email = user_details.Email
	userDetails.Phone = user_details.Phone

	tokenString, err := u.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, errors.New(errmsg.ErrCreateTocken)
	}
	return models.TokenUsers{
		Users: userDetails,
		Token: tokenString,
	}, nil
}
func (u *userUseCase) AddAddress(userID int, address models.AddressInfo) ([]models.AddressInfoResponse, error) {

	ok, err := u.helper.ValidateAlphabets(address.Name)
	if err != nil {
		return []models.AddressInfoResponse{}, errors.New(errmsg.ErrInvalidName)
	}
	if !ok {
		return []models.AddressInfoResponse{}, errors.New(errmsg.ErrInvalidName)
	}
	phone := u.helper.ValidatePhoneNumber(address.Phone)
	if !phone {
		return []models.AddressInfoResponse{}, errors.New(errmsg.ErrInvalidPhone)
	}

	pin := u.helper.ValidatePin(address.Pin)
	if !pin {
		return []models.AddressInfoResponse{}, errors.New(errmsg.ErrInvalidPin)
	}

	if userID <= 0 {
		return []models.AddressInfoResponse{}, errors.New(errmsg.ErrInvalidUId)
	}

	errs := u.userRepo.CheckUserById(userID)
	if !errs {
		return []models.AddressInfoResponse{}, errors.New(errmsg.ErrUserExist)
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
		return models.UsersProfileDetails{}, errors.New("name " + errmsg.ErrFieldEmpty)
	}
	ok, err := u.helper.ValidateAlphabets(user.Name)
	if err != nil {
		return models.UsersProfileDetails{}, errors.New(errmsg.ErrInvalidName)
	}
	if !ok {
		return models.UsersProfileDetails{}, errors.New(errmsg.ErrInvalidName)
	}
	phErr := u.helper.ValidatePhoneNumber(user.Phone)
	if !phErr {
		return models.UsersProfileDetails{}, errors.New(errmsg.ErrInvalidPhone)
	}
	details, err := u.userRepo.EditProfile(user)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return details, nil
}

func (u *userUseCase) ChangePassword(user models.ChangePassword) error {

	if user.NewPassWord == "" || user.ConfirmPassword == "" {
		return errors.New("password " + errmsg.ErrFieldEmpty)
	}

	if user.NewPassWord != user.ConfirmPassword {
		return errors.New(errmsg.ErrPasswordMatch)
	}
	newHashed, err := u.helper.PasswordHashing(user.NewPassWord)
	if err != nil {
		return errors.New(errmsg.ErrPasswordHash)
	}

	idString := strconv.FormatUint(uint64(user.UserID), 10)

	user_details, _ := u.userRepo.FindUserById(idString)

	err = u.helper.CompareHashAndPassword(user_details.Password, user.CurrentPassWord)
	if err != nil {
		return errors.New("incorrect current password")
	}

	err = u.userRepo.ChangePassword(idString, newHashed)
	if err != nil {
		return errors.New(errmsg.ErrChangePassword)
	}
	return nil
}
