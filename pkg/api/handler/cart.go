package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	CartUsecase interfaces.CartUseCase
}

func NewCartHandler(useCase interfaces.CartUseCase) *CartHandler {
	return &CartHandler{
		CartUsecase: useCase,
	}
}

func (ch *CartHandler) AddToCart(c *gin.Context) {
	var cart models.AddCart
	userID, errb := c.Get("id")
	if !errb {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided are in wrong format", nil, errors.New("getting user id is failed"))
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := c.BindJSON(&cart); err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID, _ = userID.(int)

	cartResp, err := ch.CartUsecase.AddToCart(cart)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Cannot Add to Cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	succesRsp := response.ClientResponse(http.StatusOK, "Successfully Added To Cart", cartResp, nil)
	c.JSON(http.StatusOK, succesRsp)
}

func (ch *CartHandler) ListCartItems(c *gin.Context) {
	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot list products", nil, errors.New("error in getting user id"))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	cartResp, err := ch.CartUsecase.ListCartItems(userID.(int))
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "could not get list", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	succesResp := response.ClientResponse(http.StatusOK, "successfully got the cart list", cartResp, nil)
	c.JSON(http.StatusOK, succesResp)
}

func (ch *CartHandler) UpdateProductQuantityCart(c *gin.Context) {
	var cart models.AddCart
	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot update quantity", nil, errors.New("error in getting user id"))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := c.BindJSON(&cart); err != nil {

		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	cart.UserID = userID.(int)

	cartResp, err := ch.CartUsecase.UpdateProductQuantityCart(cart)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Updation Failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	succesResp := response.ClientResponse(http.StatusOK, "Successfully Updated", cartResp, nil)
	c.JSON(http.StatusOK, succesResp)

}
