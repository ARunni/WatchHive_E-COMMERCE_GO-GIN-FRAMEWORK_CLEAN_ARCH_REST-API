package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"
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

// UserSignUp registers a new user.
// @Summary Register a new user
// @Description Registers a new user with provided details
// @Tags User
// @Accept json
// @Produce json
// @Param body body models.UserDetails true "User details for sign-up"
// @Success 201 {object} response.Response "User signed up successfully"
// @Failure 400 {object} response.Response "Invalid request or constraints not satisfied"
// @Router /user/signup [post]
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

// LoginHandler handles user login.
// @Summary Handle user login
// @Description Handles user login using provided credentials
// @Tags User
// @Accept json
// @Produce json
// @Param body body models.UserLogin true "User credentials for login"
// @Success 200 {object} response.Response "User logged in successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to log in user"
// @Router /user/login [post]
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

// AddAddress adds an address for a user.
// @Summary Add user address
// @Description Adds an address for the user identified by ID
// @Tags User Profile Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param body body models.AddressInfoResponse true "Address details for addition"
// @Success 200 {object} response.Response "Address added successfully"
// @Failure 400 {object} response.Response "Invalid request or constraints not satisfied"
// @Router /user/profile/address [post]
func (u *UserHandler) AddAddress(c *gin.Context) {
	var address models.AddressInfoResponse

	userIdstring, _ := c.Get("id")
	userId, strErr := userIdstring.(int)

	if !strErr {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, strErr)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := c.BindJSON(&address); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	err := validator.New().Struct(address)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	adrRep, err := u.userUseCase.AddAddress(userId, address)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "can not add address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusOK, "address added  successfully", adrRep, nil)
	c.JSON(http.StatusOK, successResp)

}

// ShowUserDetails retrieves details of a user.
// @Summary Retrieve user details
// @Description Retrieves details of the user identified by ID
// @Tags User Profile Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} response.Response "User details retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to get user details"
// @Router /user/profile [get]
func (u *UserHandler) ShowUserDetails(c *gin.Context) {
	userIdstring, _ := c.Get("id")
	userId, strErr := userIdstring.(int)
	if !strErr {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, strErr)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	userResp, err := u.userUseCase.ShowUserDetails(userId)
	if err != nil {
		errREsp := response.ClientResponse(http.StatusBadRequest, "Cannot get details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errREsp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "successfully got details", userResp, nil)
	c.JSON(http.StatusOK, successResp)
}

// GetAllAddress retrieves all addresses of a user.
// @Summary Retrieve all user addresses
// @Description Retrieves all addresses of the user identified by ID
// @Tags User Profile Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Success 200 {object} response.Response "All user addresses retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to get user addresses"
// @Router /user/profile/alladdress [get]
func (u *UserHandler) GetAllAddress(c *gin.Context) {
	userIdstring, _ := c.Get("id")
	userId, strErr := userIdstring.(int)
	if !strErr {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, strErr)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	userRep, err := u.userUseCase.GetAllAddress(userId)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "cannot get Addresses", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully Got All Addresses", userRep, nil)
	c.JSON(http.StatusOK, successResp)
}

// EditProfile updates user profile details.
// @Summary Update user profile
// @Description Updates user profile details based on provided information
// @Tags User Profile Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param body body models.UsersProfileDetails true "User profile details for update"
// @Success 200 {object} response.Response "User profile updated successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to update user profile"
// @Router /user/profile [patch]
func (u *UserHandler) EditProfile(c *gin.Context) {
	var details models.UsersProfileDetails

	userString, er := c.Get("id")
	if !er {
		errREsp := response.ClientResponse(http.StatusBadRequest, "Failed to get user id", nil, er)
		c.JSON(http.StatusBadRequest, errREsp)
		return
	}
	userid, ers := userString.(int)
	if !ers {
		errResp := response.ClientResponse(http.StatusBadRequest, "conversion error", nil, ers)
		c.JSON(http.StatusBadRequest, errResp)
	}
	if err := c.BindJSON(&details); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	details.ID = uint(userid)

	err := validator.New().Struct(details)
	if err != nil {
		err = errors.New("check the values that you are providing or wrong data provided")
		errResp := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	userResp, err := u.userUseCase.EditProfile(details)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot update profile", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	succesResp := response.ClientResponse(http.StatusOK, "Suceessfully updated profile", userResp, nil)
	c.JSON(http.StatusOK, succesResp)
}

// ChangePassword changes the user's password.
// @Summary Change user password
// @Description Changes the password for the user identified by ID
// @Tags User Profile Management
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Param body body models.ChangePassword true "Password details for change"
// @Success 200 {object} response.Response "Password changed successfully"
// @Failure 400 {object} response.Response "Invalid request or unable to change password"
// @Router /user/profile/password [patch]
func (u *UserHandler) ChangePassword(c *gin.Context) {
	var change models.ChangePassword

	userString, er := c.Get("id")
	if !er {

		errREsp := response.ClientResponse(http.StatusBadRequest, "Failed to get user id", nil, er)
		c.JSON(http.StatusBadRequest, errREsp)
		return
	}
	userid, ers := userString.(int)
	if !ers {
		errResp := response.ClientResponse(http.StatusBadRequest, "conversion error", nil, ers)
		c.JSON(http.StatusBadRequest, errResp)
	}
	if err := c.BindJSON(&change); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "fields provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	change.UserID = uint(userid)

	err := u.userUseCase.ChangePassword(change)

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot change password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	succesResp := response.ClientResponse(http.StatusOK, "Password changed Successfully", nil, nil)
	c.JSON(http.StatusOK, succesResp)
}
