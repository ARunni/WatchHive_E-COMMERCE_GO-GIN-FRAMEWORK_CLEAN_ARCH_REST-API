package handler

import (
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

// AddProduct adds a new product.
// @Summary Add product
// @Description Adds a new product using the provided details and image.
// @Tags Admin Product Management
// @Accept multipart/form-data
// @Produce json
// @Security BearerTokenAuth
// @Param category_id formData integer true "Category ID"
// @Param product_name formData string true "Product name"
// @Param color formData string true "Product color"
// @Param stock formData integer true "Product stock"
// @Param price formData number true "Product price"
// @Param image formData file true "Product image"
// @Success 200 {object} response.Response  "Success: Product added successfully"
// @Failure 400 {object} response.Response  "Bad request: Retrieving image error or could not add the product"
// @Router /admin/product [post]
func (i *ProductHandler) AddProduct(c *gin.Context) {
	var products models.AddProducts

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

// ListProductsUser lists products for users with pagination.
// @Summary List products for users
// @Description Retrieves a paginated list of products available for users.
// @Tags User Product Management
// @Accept json
// @Produce json
// @Param page query integer false "Page number (default: 1)"
// @Param per_page query integer false "Number of products per page (default: 5)"
// @Success 200 {object} response.Response  "Success: Products for users displayed successfully"
// @Failure 400 {object} response.Response  "Bad request: Product display error"
// @Router /user/product [get]
func (i *ProductHandler) ListProductsUser(c *gin.Context) {

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

// ListProducts lists products with pagination.
// @Summary List products
// @Description Retrieves a paginated list of products.
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param page query integer false "Page number (default: 1)"
// @Param per_page query integer false "Number of products per page (default: 5)"
// @Success 200 {object} response.Response  "Success: Products displayed successfully"
// @Failure 400 {object} response.Response  "Bad request: Product display error"
// @Router /admin/product [get]
func (i *ProductHandler) ListProductsAdmin(c *gin.Context) {

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

// EditProduct updates an existing product.
// @Summary Edit product
// @Description Updates an existing product using the provided details.
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param Product body models.ProductEdit true "Product details to be updated"
// @Success 200 {object} response.Response  "Success: Product edited successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields are in the wrong format or could not edit the product"
// @Router /admin/product [patch]
func (u *ProductHandler) EditProduct(c *gin.Context) {
	var product models.ProductEdit

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

// DeleteProduct deletes an existing product by ID.
// @Summary Delete product
// @Description Deletes an existing product by the provided ID.
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query string true "Product ID to be deleted"
// @Success 200 {object} response.Response  "Success: Product deleted successfully"
// @Failure 400 {object} response.Response  "Bad request: Product ID provided in wrong format or deletion error"
// @Router /admin/product [delete]
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

// UpdateProduct updates the stock of an existing product.
// @Summary Update product stock
// @Description Updates the stock of an existing product using the provided details.
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param ProductUpdate body models.ProductUpdate true "Product details for stock update"
// @Success 200 {object} response.Response  "Success: Product stock updated successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in wrong format or could not update the product stock"
// @Router /admin/product/stock [patch]
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
