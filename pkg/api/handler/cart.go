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

// AddToCart adds an item to the user's cart.
// @Summary Add item to cart
// @Description Adds an item to the user's cart based on the provided details.
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param AddCart body models.AddCartR true "Item details to add to the cart"
// @Success 200 {object} response.Response "Success: Item added to cart successfully"
// @Failure 400 {object} response.Response "Bad request: Fields are provided in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Cannot add item to cart"
// @Router /user/cart [post]
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

// ListCartItems retrieves the list of items in the user's cart.
// @Summary Retrieve cart items
// @Description Retrieves the list of items in the user's cart based on the user ID.
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response "Success: Retrieved cart items successfully"
// @Failure 400 {object} response.Response  "Bad request: Cannot list products"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not get the cart list"
// @Router /user/cart [get]
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

// UpdateProductQuantityCart updates the quantity of a product in the user's cart.
// @Summary Update product quantity in cart
// @Description Updates the quantity of a product in the user's cart based on the provided details.
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param UpdateCart body models.AddCartR true "Product details to update quantity"
// @Success 200 {object} response.Response  "Success: Quantity updated successfully"
// @Failure 400 {object} response.Response  "Bad request: Cannot update quantity or fields are provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Updation failed"
// @Router /user/cart [patch]
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

// RemoveFromCart removes a product from the user's cart.
// @Summary Remove product from cart
// @Description Removes a product from the user's cart based on the provided details.
// @Tags User Cart Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param RemoveFromCart body models.RemoveFromCartR true "Product details to remove from cart"
// @Success 200 {object} response.Response  "Success: Product removed from cart successfully"
// @Failure 400 {object} response.Response  "Bad request: Cannot remove product or fields are provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Removing from cart failed"
// @Router /user/cart [delete]
func (ch *CartHandler) RemoveFromCart(c *gin.Context) {
	var cart models.RemoveFromCart

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
	cartResp, err := ch.CartUsecase.RemoveFromCart(cart)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Removing from cart is Failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully Removed", cartResp, nil)
	c.JSON(http.StatusOK, successResp)

}
