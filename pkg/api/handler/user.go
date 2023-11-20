package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUseCase interfaces.UserUseCase
}

type Response struct {
	ID      uint   `copier:"must"`
	Name    string `copier:"must"`
	Surname string `copier:"must"`
}

func NewUserHandler(usecase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

func (u *UserHandler) UserSignUp(c *gin.Context) {

	var user models.UserDetails

	// bind the user details to the struct

	if err := c.BindJSON(&user); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	// checking whether the data sent by the user has all the correct constraints specified by Users struct

	err := validator.New().Struct(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	// business logic goes inside this function

	userCreated, err := u.userUseCase.UserSignUp(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "user could not signed up", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusCreated, "user signed up succesfully", userCreated, nil)
	c.JSON(http.StatusCreated, successResp)

}

func (u *UserHandler) LoginHandler(c *gin.Context) {
	var user models.UserLogin

	if err := c.BindJSON(&user); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields proveded in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	err := validator.New().Struct(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	user_details, err := u.userUseCase.LoginHandler(user)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "user could not be logged in", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "user sigend in successfully", user_details, nil)
	c.JSON(http.StatusOK, successResp)

}
