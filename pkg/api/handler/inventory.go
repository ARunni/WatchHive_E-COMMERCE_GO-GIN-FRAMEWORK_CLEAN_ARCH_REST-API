package handler

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	InventoryUseCase interfaces.InventoryUseCase
}

func NewInventoryHandler(usecase interfaces.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{
		InventoryUseCase: usecase,
	}
}

func (i *InventoryHandler) AddInventory(c *gin.Context) {
	var inventory models.AddInventories

	if err := c.ShouldBindJSON(&inventory); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	InventoryResponse, err := i.InventoryUseCase.AddInventory(inventory)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Inventory", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added inventory", InventoryResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *InventoryHandler) ListProducts(c *gin.Context) {

	pageNo := c.DefaultQuery("page", "1")
	pageList := c.DefaultQuery("per_page", "5")
	pageNoInt, err := strconv.Atoi(pageNo)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	pageListInt, err := strconv.Atoi(pageList)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
	}

	products_list, err := i.InventoryUseCase.ListProducts(pageNoInt, pageListInt)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	message := "product list"

	successRes := response.ClientResponse(http.StatusOK, message, products_list, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *InventoryHandler) EditInventory(c *gin.Context) {
	var inventory domain.Inventory

	id := c.Query("inventory_id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "problems in the id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := c.BindJSON(&inventory); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	modInventory, err := u.InventoryUseCase.EditInventory(inventory, idInt)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "could not edit the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "sucessfully edited products", modInventory, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *InventoryHandler) DeleteInventory(c *gin.Context) {

	inventoryID := c.Query("id")

	err := u.InventoryUseCase.DeleteInventory(inventoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Sucessfully deleted the product", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *InventoryHandler) UpdateInventory(c *gin.Context) {

	var p models.InventoryUpdate

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fileds are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := i.InventoryUseCase.UpdateInventory(p.Productid, p.Stock)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could  not update the inventory stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully upadated inventory stock", a, nil)
	c.JSON(http.StatusOK, successRes)
}
