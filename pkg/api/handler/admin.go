package handler

import (
	"WatchHive/pkg/helper"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"

	"net/http"
	"strconv"
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
	// c.Set("Refresh", admin.RefreshToken)

	successResp := response.ClientResponse(http.StatusOK, "logined successfully", admin, nil)
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

func (ad *AdminHandler) BlockUser(c *gin.Context) {
	id := c.Query("id")
	err := ad.adminUseCase.BlockUser(id)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	succesResp := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, succesResp)

}
func (ad *AdminHandler) UnBlockUser(c *gin.Context) {

	id := c.Query("id")
	err := ad.adminUseCase.UnBlockUser(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (ad *AdminHandler) GetUsers(c *gin.Context) {

	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ad.adminUseCase.GetUsers(page)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	c.JSON(http.StatusOK, successRes)

}

