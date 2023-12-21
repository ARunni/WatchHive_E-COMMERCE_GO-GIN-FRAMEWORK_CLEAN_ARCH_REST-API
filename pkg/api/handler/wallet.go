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
