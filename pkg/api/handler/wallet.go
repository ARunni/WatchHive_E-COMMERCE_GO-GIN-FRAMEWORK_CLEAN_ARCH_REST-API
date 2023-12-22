package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/response"
	"fmt"
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
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response  "Success: Wallet details retrieved successfully"
// @Failure 400 {object} response.Response  "Bad request: User ID not found or invalid user ID type"
// @Failure 500 {object} response.Response  "Internal server error: Failed to retrieve wallet details"
// @Router /wallet [get]
func (wh *WalletHandler) GetWallet(c *gin.Context) {
	userId, exists := c.Get("id")
	fmt.Println("userid", userId)
	if !exists {
		errs := response.ClientResponse(http.StatusBadRequest, "user_id not found", nil, "user_id is required")
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	userID, ok := userId.(int)
	if !ok {
		errs := response.ClientResponse(http.StatusBadRequest, "invalid user_id type", nil, "user_id must be an integer")
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	WalletDetails, err := wh.walletUsecase.GetWallet(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Wallet Details", WalletDetails, nil)
	c.JSON(http.StatusOK, success)
}
