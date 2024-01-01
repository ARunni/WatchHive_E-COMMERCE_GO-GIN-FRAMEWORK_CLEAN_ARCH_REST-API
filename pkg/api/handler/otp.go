package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpUseCase interfaces.OtpUseCase
}

func NewOtpHandler(usecase interfaces.OtpUseCase) *OtpHandler {
	return &OtpHandler{
		otpUseCase: usecase,
	}

}

// SendOTP sends an OTP to a provided phone number.
// @Summary Send OTP
// @Description Sends an OTP (One-Time Password) to the provided phone number for verification.
// @Tags User
// @Accept json
// @Produce json
// @Param OTPdata body models.OTPdata true "Phone number to send OTP"
// @Success 200 {object} response.Response  "Success: OTP sent successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in wrong format or OTP not sent"
// @Router /user/otplogin [post]
func (ot *OtpHandler) SendOTP(c *gin.Context) {

	var phone models.OTPdata
	if err := c.BindJSON(&phone); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}

	err := ot.otpUseCase.SendOTP(phone.PhoneNumber)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgOTPSentErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgOTPVerifySuccess, nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// VerifyOTP verifies the provided OTP code.
// @Summary Verify OTP
// @Description Verifies the provided OTP (One-Time Password) code for user authentication.
// @Tags User
// @Accept json
// @Produce json
// @Param VerifyData body models.VerifyData true "Data containing OTP for verification"
// @Success 200 {object} response.Response  "Success: OTP verified successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided are in wrong format or could not verify OTP"
// @Router /user/verifyotp [post]
func (ot *OtpHandler) VerifyOTP(c *gin.Context) {

	var code models.VerifyData
	if err := c.BindJSON(&code); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ot.otpUseCase.VerifyOTP(code)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgOTPVerifyErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgOTPVerifySuccess, users, nil)
	c.JSON(http.StatusOK, successRes)

}
