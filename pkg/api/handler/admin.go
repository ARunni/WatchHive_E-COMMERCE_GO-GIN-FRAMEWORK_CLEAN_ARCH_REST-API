package handler

import (
	"WatchHive/pkg/helper"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AdminHandler struct {
	adminUseCase interfaces.AdminUseCase
}

func NewAdminHandler(usecase interfaces.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: usecase,
	}
}

func (ad *AdminHandler) LoginHandler(c *gin.Context) {

	var adminDetails models.AdminLogin

	if err := c.BindJSON(&adminDetails); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "details is not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	admin, err := ad.adminUseCase.LoginHandler(adminDetails)
	if err != nil {
		errREsp := response.ClientResponse(http.StatusBadRequest, "cannot authenticate user", nil, err.Error())
		c.JSON(http.StatusBadRequest, errREsp)
		return
	}

	c.Set("Access", admin.AccessToken)
	c.Set("Refresh", admin.RefreshToken)

	successResp := response.ClientResponse(http.StatusOK, "logined successfully", admin, err.Error())
	c.JSON(http.StatusOK, successResp)

}

func (ad *AdminHandler) ValidateRefreshTokenAndCreateNewAccess(c *gin.Context) {
	refreshToken := c.Request.Header.Get("RefreshToken")

	// Check  refresh token is valid.

	_, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("refreshsecret"), nil
	})
	if err != nil {
		c.AbortWithError(401, errors.New("refresh token is invalid : user have to login again "))
		return
	}
	claims := &helper.AuthCustomClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	newAccessToken, err := token.SignedString([]byte("accesssecret"))
	if err != nil {
		c.AbortWithError(500, errors.New("error in creating new accesstoken"))
	}
	c.JSON(200, newAccessToken)
}
