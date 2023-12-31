package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	CouponUsecase interfaces.CouponUsecase
}

func NewCouponHandler(coupon interfaces.CouponUsecase) *CouponHandler {
	return &CouponHandler{CouponUsecase: coupon}
}

// AddCoupon Adding coupons
// @Summary Add a new coupon
// @Description Adds a new coupon based on the provided details.
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param AddCoupon body models.Coupon true "Coupon details to add"
// @Success 200 {object} response.Response "Success: Coupon added successfully"
// @Failure 400 {object} response.Response "Bad request: Fields are provided in the wrong format"
// @Failure 401 {object} response.Response "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response "Internal server error: Could not add the coupon"
// @Router /admin/coupon [post]
func (ch *CouponHandler) AddCoupon(c *gin.Context) {
	var coupon models.Coupon
	if err := c.BindJSON(&coupon); err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	couponRep, err := ch.CouponUsecase.AddCoupon(coupon)
	if err != nil {
		errResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgCouponAddFailed, nil, err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	successResp := response.ClientResponse(http.StatusOK, errmsg.MsgAddSuccess, couponRep, nil)
	c.JSON(http.StatusOK, successResp)
}

// GetCoupon Get All coupons
// @Summary Get coupons
// @Description Retrieves coupon information
// @Tags Admin Coupon Management
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success: Retrieved coupons successfully"
// @Failure 400 {object} response.Response "Bad request: Error while retrieving coupons"
// @Security BearerTokenAuth
// @Router /admin/coupon [get]
func (ch *CouponHandler) GetCouponAdmin(c *gin.Context) {
	couponResp, err := ch.CouponUsecase.GetCoupon()
	if err != nil {
		ereResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, ereResp)
	}
	successRep := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, couponResp, nil)
	c.JSON(http.StatusOK, successRep)
}

// GetCoupon Get All coupons
// @Summary Get coupons
// @Description Retrieves coupon information
// @Tags User Coupon Management
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success: Retrieved coupons successfully"
// @Failure 400 {object} response.Response "Bad request: Error while retrieving coupons"
// @Security BearerTokenAuth
// @Router /user/coupon [get]
func (ch *CouponHandler) GetCouponUser(c *gin.Context) {
	couponResp, err := ch.CouponUsecase.GetCoupon()
	if err != nil {
		ereResp := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGetErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, ereResp)
	}
	successRep := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, couponResp, nil)
	c.JSON(http.StatusOK, successRep)
}
