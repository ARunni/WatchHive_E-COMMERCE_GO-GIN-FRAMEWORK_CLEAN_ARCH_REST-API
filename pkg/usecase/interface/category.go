package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type CategoryUseCase interface {
	AddCategory(category models.CategoryAdd) (domain.Category, error)
	GetCategories() ([]domain.Category, error)
	UpdateCategory(currentId int, new string) (domain.Category, error)
	DeleteCategory(categoryID string) error
}
