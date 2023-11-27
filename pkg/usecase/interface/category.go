package interfaces

import "WatchHive/pkg/domain"

type CategoryUseCase interface {
	AddCategory(category domain.Category) (domain.Category, error)
	GetCategories() ([]domain.Category, error)
	UpdateCategory(currentId int, new string) (domain.Category, error)
	DeleteCategory(categoryID string) error
}
