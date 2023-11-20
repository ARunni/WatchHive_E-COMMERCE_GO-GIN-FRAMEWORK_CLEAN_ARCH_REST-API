package interfaces

import "WatchHive/pkg/utils/models"

type Helper interface {
	GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error)
	GenerateTokenClients(user models.UserDetailsResponse) (string,error)
	PasswordHashing(string) (string,error)
	CompareHashAndPassword(a string, b string) error
	Copy(udr *models.UserDetailsResponse, usr *models.UserSignInResponse) (models.UserDetailsResponse,error)
}
 