package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase   interfaces.OrderUseCase
	paymentUsecase interfaces.PaymentUseCase
}

func NewOrderHandler(oUsecase interfaces.OrderUseCase, pUsecase interfaces.PaymentUseCase) *OrderHandler {
	return &OrderHandler{
		orderUsecase:   oUsecase,
		paymentUsecase: pUsecase,
	}
}

func (oh *OrderHandler) CheckOut(c *gin.Context) {
	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, "Getting ID failed", nil, errors.New("failed to get id").Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	checkOutResp, err := oh.orderUsecase.Checkout(userID.(int))

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "CheckOut Failed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusOK, "Successfully completed", checkOutResp, nil)
	c.JSON(http.StatusOK, successResp)
}

func (oh *OrderHandler) OrderItemsFromCart(c *gin.Context) {

	id, errs := c.Get("id")
	if !errs {
		err := errors.New("error in getting id")
		errorRes := response.ClientResponse(http.StatusBadRequest, "bad request", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	userID := id.(int)
	var orderFromCart models.OrderFromCart
	if err := c.ShouldBindJSON(&orderFromCart); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "bad request", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orderSuccessResponse, err := oh.orderUsecase.OrderItemsFromCart(orderFromCart, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully created the order", orderSuccessResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

func (oh *OrderHandler) PlaceOrderCOD(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	paymentMethodID, err := oh.paymentUsecase.PaymentMethodID(order_id)
	if err != nil {
		err := response.ClientResponse(http.StatusInternalServerError, "error from paymentId ", nil, err.Error())
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if paymentMethodID == 1 {
		err := oh.orderUsecase.ExecutePurchaseCOD(order_id)
		if err != nil {
			errorRes := response.ClientResponse(http.StatusInternalServerError, "error in cash on delivery ", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorRes)
			return
		}
		success := response.ClientResponse(http.StatusOK, "Placed Order with cash on delivery", nil, nil)
		c.JSON(http.StatusOK, success)
	}
	if paymentMethodID == 2 {
		success := response.ClientResponse(http.StatusOK, "Placed Order with razor pay", nil, nil)
		c.JSON(http.StatusOK, success)
	}
}

func (oh *OrderHandler) GetOrderDetails(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	id, errs := c.Get("id")

	if !errs {
		if err != nil {
			err := errors.New("couldn't get id ")
			errorRes := response.ClientResponse(http.StatusInternalServerError, "Error in getting id", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorRes)
			return
		}
	}

	UserID := id.(int)
	OrderDetails, err := oh.orderUsecase.GetOrderDetails(UserID, page, pageSize)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Full Order Details", OrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}
