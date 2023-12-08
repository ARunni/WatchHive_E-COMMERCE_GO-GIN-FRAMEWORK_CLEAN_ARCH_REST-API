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

type ProductHandler struct {
	ProductUseCase interfaces.ProductUseCase
}

func NewProductHandler(usecase interfaces.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: usecase,
	}
}

func (i *ProductHandler) AddProduct(c *gin.Context) {
	var products models.AddProducts

	// if err := c.ShouldBindJSON(&product); err != nil {
	// 	errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
	// 	c.JSON(http.StatusBadRequest, errorRes)
	// 	return
	// }

	cat := c.PostForm("category_id")
	products.CategoryID, _ = strconv.Atoi(cat)
	products.ProductName = c.PostForm("product_name")
	products.Color = c.PostForm("color")
	products.Stock, _ = strconv.Atoi(c.PostForm("stock"))
	products.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

	file, err := c.FormFile("image")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from the Form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	ProductResponse, err := i.ProductUseCase.AddProduct(products, file)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added Product", ProductResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *ProductHandler) ListProducts(c *gin.Context) {

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

	products_list, err := i.ProductUseCase.ListProducts(pageNoInt, pageListInt)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Product cannot be displayed", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	message := "product list"

	successRes := response.ClientResponse(http.StatusOK, message, products_list, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *ProductHandler) EditProduct(c *gin.Context) {
	var product domain.Product

	// id := c.Query("product_id")
	// idInt, err := strconv.Atoi(id)

	// if err != nil {
	// 	errRes := response.ClientResponse(http.StatusBadRequest, "problems in the id", nil, err.Error())
	// 	c.JSON(http.StatusBadRequest, errRes)
	// 	return
	// }

	if err := c.BindJSON(&product); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	modProduct, err := u.ProductUseCase.EditProduct(product)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "could not edit the product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "sucessfully edited products", modProduct, nil)
	c.JSON(http.StatusOK, successRes)
}

func (u *ProductHandler) DeleteProduct(c *gin.Context) {

	productID := c.Query("id")

	err := u.ProductUseCase.DeleteProduct(productID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Sucessfully deleted the product", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *ProductHandler) UpdateProduct(c *gin.Context) {

	var p models.ProductUpdate

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fileds are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := i.ProductUseCase.UpdateProduct(p.Productid, p.Stock)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could  not update the product stock", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully upadated product stock", a, nil)
	c.JSON(http.StatusOK, successRes)
}
