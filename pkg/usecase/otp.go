package usecase

import (
	"WatchHive/pkg/config"
	helper "WatchHive/pkg/helper/interface"
	repo "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"

	"github.com/jinzhu/copier"
)

type otpUseCase struct {
	cfg           config.Config
	otpRepository repo.OtpRepository
	helper        helper.Helper
}

func NewOtpUseCase(cfg config.Config, repo repo.OtpRepository, h helper.Helper) interfaces.OtpUseCase {
	return &otpUseCase{
		cfg:           cfg,
		otpRepository: repo,
		helper:        h,
	}

}

func (ot *otpUseCase) SendOTP(phone string) error {

	phoneNumber := ot.helper.ValidatePhoneNumber(phone)

	if !phoneNumber {
		return errors.New(errmsg.ErrInvalidPhone)
	}

	ok := ot.otpRepository.FindUserByMobileNumber(phone)

	if !ok {
		return errors.New(errmsg.ErrUserExistFalse)
	}

	ot.helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	_, err := ot.helper.TwilioSendOTP(phone, ot.cfg.SERVICESID)

	if err != nil {
		return err
	}
	return nil

}

func (ot *otpUseCase) VerifyOTP(code models.VerifyData) (models.TokenUsers, error) {

	phoneNumber := ot.helper.ValidatePhoneNumber(code.PhoneNumber)

	if !phoneNumber {
		return models.TokenUsers{}, errors.New(errmsg.ErrInvalidPhone)
	}

	ot.helper.TwilioSetup(ot.cfg.ACCOUNTSID, ot.cfg.AUTHTOKEN)
	err := ot.helper.TwilioVerifyOTP(ot.cfg.SERVICESID, code.Code, code.PhoneNumber)

	if err != nil {
		//this guard clause catches the error code runs only until here
		return models.TokenUsers{}, errors.New(errmsg.ErrVerify)
	}

	// if user is authenticated using OTP send back user details

	userDetails, err := ot.otpRepository.UserDetailsUsingPhone(code.PhoneNumber)
	if err != nil {
		return models.TokenUsers{}, err
	}

	tokenString, err := ot.helper.GenerateTokenClients(userDetails)

	if err != nil {
		return models.TokenUsers{}, err
	}

	var user models.UserDetailsResponse
	err = copier.Copy(&user, &userDetails)

	if err != nil {
		return models.TokenUsers{}, err
	}

	return models.TokenUsers{
		Users: user,
		Token: tokenString,
	}, nil

}
