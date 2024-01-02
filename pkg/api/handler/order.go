package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"errors"
	"fmt"
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

// CheckOut processes the checkout for the user's order.
// @Summary Process checkout
// @Description Processes the checkout for the user's order.
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response  "Success: Checkout completed successfully"
// @Failure 400 {object} response.Response  "Bad request: Getting user ID failed"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Checkout failed"
// @Router /user/orders/checkout [get]
func (oh *OrderHandler) CheckOut(c *gin.Context) {
	userID, errs := c.Get("id")
	if !errs {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr+"ID", nil, errors.New("failed to get id").Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	checkOutResp, err := oh.orderUsecase.Checkout(userID.(int))

	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgCheckoutErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, checkOutResp, nil)
	c.JSON(http.StatusOK, successResp)
}

// OrderItemsFromCart places an order with items from the user's cart.
// @Summary Place order from cart
// @Description Places an order with items from the user's cart based on the provided details.
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param OrderFromCart body models.OrderFromCart true "Order details from cart"
// @Success 200 {object} response.Response  "Success: Order placed successfully"
// @Failure 400 {object} response.Response  "Bad request: Error in getting ID or bad request"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not place the order"
// @Router /user/orders [post]
func (oh *OrderHandler) OrderItemsFromCart(c *gin.Context) {

	id, errs := c.Get("id")
	if !errs {
		err := errors.New("error in getting id")
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgBadRequestErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	userID := id.(int)
	var orderFromCart models.OrderFromCart
	if err := c.ShouldBindJSON(&orderFromCart); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgBadRequestErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orderSuccessResponse, err := oh.orderUsecase.OrderItemsFromCart(orderFromCart, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgAddErr+"order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	success := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, orderSuccessResponse, nil)
	c.JSON(http.StatusOK, success)
}

// GetOrderDetails retrieves order details for a user.
// @Summary Retrieve order details
// @Description Retrieves order details for a user based on the provided pagination parameters.
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param count query integer false "Number of items per page (default: 10)"
// @Success 200 {object} response.Response  "Success: Retrieved order details successfully"
// @Failure 400 {object} response.Response  "Bad request: Page number or count not in correct format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not retrieve order details"
// @Router /user/orders [get]
func (oh *OrderHandler) GetOrderDetails(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgPageNumFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgPageCountErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	id, errs := c.Get("id")

	if !errs {
		if err != nil {
			err := errors.New("couldn't get id ")
			errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr+"id", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorRes)
			return
		}
	}

	UserID := id.(int)
	OrderDetails, err := oh.orderUsecase.GetOrderDetails(UserID, page, pageSize)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgOrderErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, OrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// CancelOrder cancels an order by ID for the logged-in user.
// @Summary Cancel order
// @Description Cancels an order based on the provided order ID for the logged-in user.
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query integer true "Order ID to cancel"
// @Success 200 {object} response.Response  "Success: Order canceled successfully"
// @Failure 400 {object} response.Response  "Bad request: Error from orderID or error from userid"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not cancel the order"
// @Router /user/orders [delete]
func (oh *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgErr+"from orderID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	id, errs := c.Get("id")
	if !errs {
		err := errors.New("error in getting id")
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgErr+"from userid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userID := id.(int)
	err = oh.orderUsecase.CancelOrders(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgCanelErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, nil, nil)
	c.JSON(http.StatusOK, success)
}

// GetAllOrderDetailsForAdmin retrieves all order details for admin with pagination.
// @Summary Retrieve all order details for admin
// @Description Retrieves all order details for admin with pagination based on the provided parameters.
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param size query integer false "Number of items per page (default: 10)"
// @Success 200 {object} response.Response  "Success: Retrieved all order details for admin successfully"
// @Failure 400 {object} response.Response  "Bad request: Page number or count not in correct format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not retrieve order details for admin"
// @Router /admin/orders [get]
func (oh *OrderHandler) GetAllOrderDetailsForAdmin(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgPageNumFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	countStr := c.DefaultQuery("size", "10")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgPageNumFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	var pageStruct models.Page
	pageStruct.Page = page
	pageStruct.Size = pageSize
	allOrderDetails, err := oh.orderUsecase.GetAllOrdersAdmin(pageStruct)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, allOrderDetails, nil)
	c.JSON(http.StatusOK, success)
}

// ApproveOrder approves an order by its ID.
// @Summary Approve order
// @Description Approves an order based on the provided order ID.
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param order_id query integer true "Order ID to approve"
// @Success 200 {object} response.Response  "Success: Order approved successfully"
// @Failure 400 {object} response.Response  "Bad request: Error from orderID or couldn't approve the order"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Failed to approve the order"
// @Router /admin/orders [patch]
func (oh *OrderHandler) ApproveOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgErr+"from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = oh.orderUsecase.ApproveOrder(orderId)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgOrderApproveErr, nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, nil, nil)
	c.JSON(http.StatusOK, success)
}

// CancelOrderFromAdmin cancels an order by its ID from an admin perspective.
// @Summary Cancel order from admin
// @Description Cancels an order based on the provided order ID from an admin perspective.
// @Tags Admin Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param order_id query integer true "Order ID to cancel"
// @Success 200 {object} response.Response  "Success: Order canceled successfully"
// @Failure 500 {object} response.Response  "Internal server error: Error from orderID or couldn't cancel the order"
// @Router /admin/orders [delete]
func (oh *OrderHandler) CancelOrderFromAdmin(c *gin.Context) {
	order_id, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, errmsg.MsgErr+"from orderID", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	err = oh.orderUsecase.CancelOrderFromAdmin(order_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, errmsg.MsgCanelErr, nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, nil, nil)
	c.JSON(http.StatusOK, success)
}

// ReturnOrder initiates the return process for a specific order.
// @Summary Initiate order return
// @Description Initiates the return process for an order based on the provided order ID and user ID.
// @Tags  User Order Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param order_id query integer true "Order ID to initiate return"
// @Success 200 {object} response.Response  "Success: Order returned successfully"
// @Failure 400 {object} response.Response  "Bad request: Error from orderID or error from userid"
// @Failure 500 {object} response.Response  "Internal server error: Couldn't initiate the order return"
// @Router /user/orders [patch]
func (oh *OrderHandler) ReturnOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Query("order_id"))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgErr+"from orderID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userId, errs := c.Get("id")
	if !errs {
		err := errors.New("error in getting id")
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgErr+"from userID", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userID := userId.(int)
	err = oh.orderUsecase.ReturnOrder(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, errmsg.MsgCanelErr, nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, nil, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary Print Invoice
// @Description Generate and download the invoice for a specific order.
// @Accept json
// @Produce application/pdf
// @Tags User Order Management
// @security BearerTokenAuth
// @Param order_id query string true "Order ID for which the invoice should be generated"
// @Success 200 {file} pdf "Invoice PDF"
// @Failure 400 {object} response.Response "Error in processing the request"
// @Failure 500 {object} response.Response "Error generating or downloading the invoice"
// @Router /user/orders/print [get]
func (O *OrderHandler) PrintInvoice(c *gin.Context) {
	userId, _ := c.Get("id")
	userID := userId.(int)

	orderId := c.Query("order_id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		err = errors.New(errmsg.ErrDatatypeConversion + err.Error())
		errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgIdErr, nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	pdf, err := O.orderUsecase.PrintInvoice(orderIdInt, userID)
	fmt.Println("error ", err)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgPrintErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Header("Content-Disposition", "attachment;filename=invoice.pdf")

	pdfFilePath := "salesReport/invoice.pdf"

	err = pdf.OutputFileAndClose(pdfFilePath)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgPrintErr, nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")
	c.Header("Content-Type", "application/pdf")

	c.File(pdfFilePath)

	c.Header("Content-Type", "application/pdf")

	err = pdf.Output(c.Writer)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, errmsg.MsgPrintErr, nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, pdf, nil)
	c.JSON(http.StatusOK, successRes)
}
