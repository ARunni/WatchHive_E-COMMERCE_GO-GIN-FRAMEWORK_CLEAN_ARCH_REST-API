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

type PaymentHandler struct {
	paymentUseCase interfaces.PaymentUseCase
}

func NewPaymentHandler(usecase interfaces.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: usecase,
	}

}
// AddPaymentMethod adds a new payment method.
// @Summary Add payment method
// @Description Adds a new payment method using the provided details.
// @Tags Admin Payment Methods
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param NewPaymentMethod body models.NewPaymentMethod true "Details of the new payment method"
// @Success 200 {object} response.Response  "Success: Payment method added successfully"
// @Failure 400 {object} response.Response  "Bad request: Cannot add payment method or payment name error"
// @Router /admin/payment [post]
func (ph *PaymentHandler) AddPaymentMethod(c *gin.Context) {
	var payment models.NewPaymentMethod

	err := c.BindJSON(&payment)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot Add payment method Payment name", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	paymentResp, err := ph.paymentUseCase.AddPaymentMethod(payment)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, "Cannot Add payment method", nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, "Successfully Added", paymentResp, nil)
	c.JSON(http.StatusOK, successResp)
}



func (ph *PaymentHandler) MakePaymentRazorpay(c *gin.Context) {

	userId := c.Query("user_id")

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		err := errors.New("error in converting string to int userid")
		errRes := response.ClientResponse(http.StatusBadRequest, "Cannot make payment", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	// userId := 5
	orderId := c.Query("order_id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		err := errors.New("error in converting string to int orderid")
		errRes := response.ClientResponse(http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	body, razorId, err := ph.paymentUseCase.MakePaymentRazorpay(orderIdInt, userIdInt)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// razorPay:=

	c.HTML(http.StatusOK, "index.html", gin.H{
		"final_price": body.FinalPrice * 100,
		"razor_id":    razorId,
		"user_id":     userId,
		"order_id":    body.OrderId,
		"user_name":   body.Name,
		"total":       int(body.FinalPrice),
		// "razor_key": razorPay,
	})
}

func (pu *PaymentHandler) VerifyPayment(c *gin.Context) {
	orderId := c.Query("order_id")
	paymentId := c.Query("payment_id")
	razorId := c.Query("razor_id")

	if err := pu.paymentUseCase.SavePaymentDetails(paymentId, razorId, orderId); err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "could not update payment details", nil, err.Error())
		c.JSON(http.StatusOK, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated payment details", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
