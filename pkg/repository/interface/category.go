package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type CategoryRepository interface {
	AddCategory(category models.CategoryAdd) (domain.Category, error)
	GetCategories() ([]domain.Category, error)
	UpdateCategory(currentId int, new string) (domain.Category, error)
	CheckCategory(currentId int) (bool, error)
	DeleteCategory(categoryID string) error
	CheckCategoryByName(name string) bool
	GetCategoryId(productId int) (int, error)
}
