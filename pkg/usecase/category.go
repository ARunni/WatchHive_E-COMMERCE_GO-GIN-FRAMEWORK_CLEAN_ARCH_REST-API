package usecase

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"
	"errors"
)

type categoryUseCase struct {
	repository interfaces.CategoryRepository
}

func NewCategoryUseCase(repo interfaces.CategoryRepository) services.CategoryUseCase {
	return &categoryUseCase{
		repository: repo,
	}

}
func (Cat *categoryUseCase) AddCategory(category domain.Category) (domain.Category, error) {

	productResponse, err := Cat.repository.AddCategory(category)

	if err != nil {
		return domain.Category{}, err
	}

	return productResponse, nil

}

func (Cat *categoryUseCase) GetCategories() ([]domain.Category, error) {
	categories, err := Cat.repository.GetCategories()
	if err != nil {
		return []domain.Category{}, err
	}
	return categories, nil
}

func (Cat *categoryUseCase) UpdateCategory(current string, new string) (domain.Category, error) {

	result, err := Cat.repository.CheckCategory(current)
	if err != nil {
		return domain.Category{}, err
	}

	if !result {
		return domain.Category{}, errors.New("there is no category as you mentioned")
	}

	newcat, err := Cat.repository.UpdateCategory(current, new)
	if err != nil {
		return domain.Category{}, err
	}

	return newcat, err
}

func (Cat *categoryUseCase) DeleteCategory(categoryID string) error {

	err := Cat.repository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil
}
