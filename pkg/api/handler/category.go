package handler

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"WatchHive/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryUseCase interfaces.CategoryUseCase
}

func NewCategoryHandler(usecase interfaces.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

func (cat *CategoryHandler) AddCategory(c *gin.Context) {

	var category domain.Category
	if err := c.BindJSON(&category); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	CategoryResponse, err := cat.CategoryUseCase.AddCategory(category)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Could not add the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully added Category", CategoryResponse, nil)
	c.JSON(http.StatusOK, successRes)
}

func (Cat *CategoryHandler) GetCategory(c *gin.Context) {

	categories, err := Cat.CategoryUseCase.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

func (Cat *CategoryHandler) UpdateCategory(c *gin.Context) {
	var p models.SetNewName

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := Cat.CategoryUseCase.UpdateCategory(p.Current, p.New)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Sucessfully updated...", a, nil)
	c.JSON(http.StatusOK, successRes)
}

func (Cat *CategoryHandler) DeleteCategory(c *gin.Context) {

	categoryID := c.Query("id")
	err := Cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Fields are not provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	SuccessRes := response.ClientResponse(http.StatusOK, "Sucessfully Deleted...", nil, nil)
	c.JSON(http.StatusOK, SuccessRes)
}
