package interfaces

import "WatchHive/pkg/domain"

type CategoryRepository interface {
	AddCategory(category domain.Category) (domain.Category,error)
	GetCategories() ([]domain.Category,error)
	UpdateCategory(current, new string) (domain.Category,error)
	CheckCategory(current string) (bool,error)
	DeleteCategory(categoryID string) error
}