package handler

import (
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
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

// AddCategory adds a new category.
// @Summary Add a new category
// @Description Adds a new category based on the provided details.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param AddCategory body models.CategoryAdd true "Category details to add"
// @Success 200 {object} response.Response  "Success: Category added successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields are provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not add the category"
// @Router /admin/category [post]
func (cat *CategoryHandler) AddCategory(c *gin.Context) {

	var category models.CategoryAdd
	if err := c.BindJSON(&category); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	CategoryResponse, err := cat.CategoryUseCase.AddCategory(category)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgAddErr+"category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, CategoryResponse, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetCategory retrieves all categories.
// @Summary Retrieve all categories
// @Description Retrieves all categories available.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Success 200 {object} response.Response  "Success: Retrieved all categories successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not retrieve categories"
// @Router /admin/category [get]
func (Cat *CategoryHandler) GetCategory(c *gin.Context) {

	categories, err := Cat.CategoryUseCase.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgGetSucces, categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// UpdateCategory updates an existing category's name.
// @Summary Update category name
// @Description Updates the name of an existing category based on the provided details.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param SetNewName body models.SetNewName true "Current and New category name details"
// @Success 200 {object} response.Response  "Success: Category name updated successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields provided in the wrong format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not update the category name"
// @Router /admin/category [patch]
func (Cat *CategoryHandler) UpdateCategory(c *gin.Context) {
	var p models.SetNewName

	if err := c.BindJSON(&p); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	a, err := Cat.CategoryUseCase.UpdateCategory(p.CurrentId, p.New)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgUpdateErr+"category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, a, nil)
	c.JSON(http.StatusOK, successRes)
}

// DeleteCategory deletes a category by ID.
// @Summary Delete category
// @Description Deletes a category based on the provided category ID.
// @Tags Admin Category Management
// @Accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id query string true "Category ID to delete"
// @Success 200 {object} response.Response  "Success: Category deleted successfully"
// @Failure 400 {object} response.Response  "Bad request: Fields are not provided in the correct format"
// @Failure 401 {object} response.Response  "Unauthorized: Invalid or missing authentication"
// @Failure 500 {object} response.Response  "Internal server error: Could not delete the category"
// @Router /admin/category [delete]
func (Cat *CategoryHandler) DeleteCategory(c *gin.Context) {

	categoryID := c.Query("id")
	err := Cat.CategoryUseCase.DeleteCategory(categoryID)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, errmsg.MsgFormatErr, nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	SuccessRes := response.ClientResponse(http.StatusOK, errmsg.MsgSuccess, nil, nil)
	c.JSON(http.StatusOK, SuccessRes)
}
