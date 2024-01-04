package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletUsecase interfaces.WalletUsecase
}

func NewWalletHandler(usecase interfaces.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: usecase}
}

// GetWallet retrieves wallet details for a user.
// @Summary Retrieve wallet details
// @Description Retrieves wallet details for a specific user.
// @Tags User Wallet Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response  "Success: Wallet details retrieved successfully"
// @Failure 400 {object} response.Response  "Bad request: User ID not found or invalid user ID type"
// @Failure 500 {object} response.Response  "Internal server error: Failed to retrieve wallet details"
// @Router /user/wallet [get]
func (wh *WalletHandler) GetWallet(c *gin.Context) {
	userId, exists := c.Get("id")

	if !exists {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgUserIdErr, nil, errmsg.MsgRequiredUserIdErr)
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	userID, ok := userId.(int)
	if !ok {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgInvalidIdErr, nil, errmsg.MsgIdDatatypeErr)
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	WalletDetails, err := wh.walletUsecase.GetWallet(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, errmsg.MsgGettingDataErr, nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, WalletDetails, nil)
	c.JSON(http.StatusOK, success)
}

// GetWalletHistory
// @Summary Get wallet history
// @Description Retrieves wallet history information for a user
// @Tags User Wallet Management
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Success: Retrieved wallet history successfully"
// @Failure 400 {object} response.Response "Bad request: Error while retrieving wallet history"
// @Security BearerTokenAuth
// @Router /user/wallet/history [get]
func (wh *WalletHandler) GetWalletHistory(c *gin.Context) {
	userId, _ := c.Get("id")
	userID, _ := userId.(int)
	walletHistory, err := wh.walletUsecase.GetWalletHistory(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, errmsg.MsgGettingDataErr, nil, err)
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	succesResp := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, walletHistory, nil)
	c.JSON(http.StatusOK, succesResp)
}
