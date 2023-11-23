package helper

import (
	"WatchHive/pkg/config"
	interfaces "WatchHive/pkg/helper/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
	"golang.org/x/crypto/bcrypt"
)

type helper struct {
	cfg config.Config
}

func NewHelper(config config.Config) interfaces.Helper {
	return &helper{
		cfg: config,
	}
}

type AuthCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func (helper *helper) GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error) {

	accesTokenClaims := &AuthCustomClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCustomClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accesToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accesTokenClaims)
	accessTokenString, err := accesToken.SignedString([]byte("accesssecret"))
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("refreshsecret"))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil

}
func (h *helper) GenerateTokenClients(user models.UserDetailsResponse) (string, error) {
	claims := &AuthCustomClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  "client",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("watchhive@123"))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *helper) PasswordHashing(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}
	hash := string(hashedPassword)
	return hash, nil
}

func (h *helper) CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}
	return nil
}

func (h *helper) Copy(udr *models.UserDetailsResponse, usr *models.UserSignInResponse) (models.UserDetailsResponse, error) {
	err := copier.Copy(udr, usr)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}
	return *udr, nil
}

var client *twilio.RestClient

func (h *helper) TwilioSetup(username string, password string) {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})

}

func (h *helper) TwilioSendOTP(phone string, serviceID string) (string, error) {
	fmt.Println("ghhkkk", phone)
	to := "+91" + phone
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {
		return " ", err
	}
	return *resp.Sid, nil

}

func (h *helper) TwilioVerifyOTP(serviceID string, code string, phone string) error {
	to := "+91" + phone
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(serviceID, params)
	if err != nil {
		return err
	}

	if *resp.Status == "approved" {
		return nil
	}

	return errors.New("failed to validate otp")

}
