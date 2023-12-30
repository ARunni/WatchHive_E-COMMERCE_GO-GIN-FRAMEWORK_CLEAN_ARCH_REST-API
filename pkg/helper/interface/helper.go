package interfaces

import (
	"WatchHive/pkg/utils/models"
	"mime/multipart"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Helper interface {
	GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error)
	GenerateTokenClients(user models.UserDetailsResponse) (string, error)
	PasswordHashing(string) (string, error)
	CompareHashAndPassword(a string, b string) error
	Copy(udr *models.UserDetailsResponse, usr *models.UserSignInResponse) (models.UserDetailsResponse, error)

	TwilioSetup(username string, password string)
	TwilioSendOTP(phone string, serviceID string) (string, error)
	TwilioVerifyOTP(serviceID string, code string, phone string) error

	ValidatePhoneNumber(phone string) bool
	ValidatePin(pin string) bool
	ValidateDatatype(data, intOrString string) (bool, error)

	AddImageToAwsS3(file *multipart.FileHeader) (string, error)

	GetTimeFromPeriod(timePeriod string) (time.Time, time.Time)
	ValidateDate(dateString string) bool

	ValidateAlphabets(data string) (bool, error)
	ConvertToExel(sales []models.OrderDetailsAdmin) (*excelize.File, error)
}
