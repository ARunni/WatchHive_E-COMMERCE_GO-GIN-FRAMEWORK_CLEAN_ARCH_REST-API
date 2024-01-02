package usecase

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	services "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
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
func (cu *categoryUseCase) AddCategory(category models.CategoryAdd) (domain.Category, error) {
	if category.Category == "" {
		return domain.Category{}, errors.New(errmsg.ErrInvalidFormat)
	}
	ok := cu.repository.CheckCategoryByName(category.Category)
	if ok {
		return domain.Category{}, errors.New(errmsg.ErrExistTrue)
	}

	productResponse, err := cu.repository.AddCategory(category)

	if err != nil {
		return domain.Category{}, err
	}

	return productResponse, nil

}

func (cu *categoryUseCase) GetCategories() ([]domain.Category, error) {

	categories, err := cu.repository.GetCategories()
	if err != nil {
		return []domain.Category{}, err
	}
	return categories, nil
}

func (cu *categoryUseCase) UpdateCategory(currentId int, new string) (domain.Category, error) {

	if currentId <= 0 {
		return domain.Category{}, errors.New(errmsg.ErrInvalidCId)
	}

	ok, err := cu.repository.CheckCategory(currentId)
	if err != nil {
		return domain.Category{}, err
	}

	if !ok {
		return domain.Category{}, errors.New(errmsg.ErrCatExistFalse)
	}

	newcat, err := cu.repository.UpdateCategory(currentId, new)
	if err != nil {
		return domain.Category{}, err
	}

	return newcat, err
}

func (cu *categoryUseCase) DeleteCategory(categoryID string) error {
	if categoryID == "" {
		return errors.New(errmsg.ErrInvalidData)
	}
	catId, catErr := strconv.Atoi(categoryID)

	if catErr != nil || catId <= 0 {
		return errors.New(errmsg.ErrInvalidCId)
	}
	ok, err := cu.repository.CheckCategory(catId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(errmsg.ErrCatExistFalse)
	}
	err = cu.repository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil
}
