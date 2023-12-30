package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OfferHandler struct {
	OfferUsecase interfaces.OfferUsecase
}

func NewOfferHandler(usecase interfaces.OfferUsecase) *OfferHandler {
	return &OfferHandler{
		OfferUsecase: usecase,
	}
}

// @Summary Add Product Offer
// @Description Add a new product offer.
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Tags Admin Offer Management
// @Param productOffer body models.ProductOfferResp true "Product offer details in JSON format"
// @Success 201 {object} response.Response "Successfully added offer"
// @Failure 400 {object} response.Response "Invalid request format or constraints not satisfied"
// @Failure 500 {object} response.Response "Failed to add offer"
// @Router /admin/offer/product-offer [post]
func (of *OfferHandler) AddProductOffer(c *gin.Context) {

	var productOffer models.ProductOfferResp

	if err := c.ShouldBindJSON(&productOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(productOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = of.OfferUsecase.AddProductOffer(productOffer)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}

// @Summary Add Category Offer
// @Description Add a new category offer.
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Tags Admin Offer Management
// @Param categoryOffer body models.CategorytOfferResp true "Category offer details in JSON format"
// @Success 201 {object} response.Response "Successfully added offer"
// @Failure 400 {object} response.Response "Invalid request format or constraints not satisfied"
// @Failure 500 {object} response.Response "Failed to add offer"
// @Router /admin/offer/category-offer [post]
func (of *OfferHandler) AddCategoryOffer(c *gin.Context) {

	var categoryOffer models.CategorytOfferResp

	if err := c.ShouldBindJSON(&categoryOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(categoryOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = of.OfferUsecase.AddCategoryOffer(categoryOffer)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}

// @Summary Get Product Offer
// @Description Retrieve all product offers.
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Tags Admin Offer Management
// @Success 200 {object} response.Response "Successfully got all offers"
// @Failure 400 {object} response.Response "Invalid request format or fields provided in the wrong format"
// @Failure 500 {object} response.Response "Failed to retrieve offers"
// @Router /admin/offer/get-product-offer [get]
func (of *OfferHandler) GetProductOffer(c *gin.Context) {

	products, err := of.OfferUsecase.GetProductOffer()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", products, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Get Category Offer
// @Description Retrieve all category offers.
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Tags Admin Offer Management
// @Success 200 {object} response.Response "Successfully got all offers"
// @Failure 400 {object} response.Response "Invalid request format or fields provided in the wrong format"
// @Failure 500 {object} response.Response "Failed to retrieve offers"
// @Router /admin/offer/get-category-offer [get]
func (of *OfferHandler) GetCategoryOffer(c *gin.Context) {

	categories, err := of.OfferUsecase.GetCategoryOffer()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Expire Product Offer
// @Description Expire a product offer by providing its ID.
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Tags Admin Offer Management
// @Param id query int true "ID of the product offer to expire"
// @Success 200 {object} response.Response "Successfully made product offer invalid"
// @Failure 400 {object} response.Response "Invalid request format or fields provided in the wrong format"
// @Failure 500 {object} response.Response "Failed to expire product offer"
// @Router /admin/offer/expire-product-offer [delete]
func (of *OfferHandler) ExpireProductOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := of.OfferUsecase.ExpireProductOffer(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made product offer invalid", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Expire Category Offer
// @Description Expire a category offer by providing its ID.
// @Accept json
// @Produce json
// @security BearerTokenAuth
// @Tags Admin Offer Management
// @Param id query int true "ID of the category offer to expire"
// @Success 200 {object} response.Response "Successfully made category offer invalid"
// @Failure 400 {object} response.Response "Invalid request format or fields provided in the wrong format"
// @Failure 500 {object} response.Response "Failed to expire category offer"
// @Router /admin/offer/expire-category-offer [delete]
func (of *OfferHandler) ExpireCategoryOffer(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := of.OfferUsecase.ExpireCategoryOffer(id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Coupon cannot be made invalid", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made category offer invalid", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
