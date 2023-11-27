package usecase

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"
	"errors"
	"strconv"
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
	if category.Category == "" {
		return domain.Category{}, errors.New("invalid category")
	}

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

func (Cat *categoryUseCase) UpdateCategory(currentId int, new string) (domain.Category, error) {

	if currentId <= 0 {
		return domain.Category{}, errors.New("invalid category id")
	}

	result, err := Cat.repository.CheckCategory(currentId)
	if err != nil {
		return domain.Category{}, err
	}

	if !result {
		return domain.Category{}, errors.New("there is no category as you mentioned")
	}

	newcat, err := Cat.repository.UpdateCategory(currentId, new)
	if err != nil {
		return domain.Category{}, err
	}

	return newcat, err
}

func (Cat *categoryUseCase) DeleteCategory(categoryID string) error {
	catId, catErr := strconv.Atoi(categoryID)

	if catErr != nil || catId <= 0 {
		return errors.New("invalid category id")
	}

	err := Cat.repository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil
}
