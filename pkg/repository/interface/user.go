package interfaces

import "WatchHive/pkg/utils/models"

type UserRepository interface {
	UserSignup(user models.UserDetails) (models.UserDetailsResponse, error)
	CheckUserAvilability(email string) bool
	FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error)
	UserBlockStatus(email string) (bool, error)
	AddAddress(userID int, address models.AddressInfo) error
	CheckIfUserAddress(userID int) bool
	CheckUserById(userID int) bool
	ShowUserDetails(userID int) (models.UsersProfileDetails, error)
	GetAllAddress(userID int) ([]models.AddressInfoResponse, error)
	EditProfile(user models.UsersProfileDetails) (models.UsersProfileDetails, error)
	ChangePassword(userID, password string) error
	FindUserById(id string) (models.UserSignInResponse, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	NewReferralEntry(userId int, userReferral string) error 
	GetUserIdFromReferralCode(ReferralCode string) (int, error) 
	UpdateReferralAmount(referralAmount float64, referredUserId , currentUserID int) error 
	AmountInRefferals(userID int) (float64, error) 
}
