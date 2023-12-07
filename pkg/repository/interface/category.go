package interfaces

import "WatchHive/pkg/domain"

type CategoryRepository interface {
	AddCategory(category domain.Category) (domain.Category,error)
	GetCategories() ([]domain.Category,error)
	UpdateCategory(currentId int , new string) (domain.Category,error)
	CheckCategory(currentId int) (bool,error)
	DeleteCategory(categoryID string) error
	CheckCategoryByName(name string) bool
}