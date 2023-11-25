package usecase

import (
	"WatchHive/pkg/domain"
	helper "WatchHive/pkg/helper/interface"
	rep "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
)

type productUseCase struct {
	repository rep.ProductRepository
	helper     helper.Helper
}

func NewProductUseCase(repo rep.ProductRepository, h helper.Helper) interfaces.ProductUseCase {
	return &productUseCase{
		repository: repo,
		helper:     h,
	}

}

func (i *productUseCase) AddProduct(inventory models.AddProducts) (models.ProductResponse, error) {

	ProductResponse, err := i.repository.AddProduct(inventory)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return ProductResponse, nil

}
func (i *productUseCase) ListProducts(pageNo, pageList int) ([]models.ProductUserResponse, error) {

	offset := (pageNo - 1) * pageList
	productList, err := i.repository.ListProducts(pageList, offset)
	if err != nil {
		return []models.ProductUserResponse{}, err
	}
	return productList, nil
}

func (usecase *productUseCase) EditProduct(product domain.Product, id int) (domain.Product, error) {
	modProduct, err := usecase.repository.EditProduct(product, id)
	if err != nil {
		return domain.Product{}, err
	}
	return modProduct, nil
}

func (usecase *productUseCase) DeleteProduct(productID string) error {

	err := usecase.repository.DeleteProduct(productID)
	if err != nil {
		return err
	}
	return nil
}

func (i productUseCase) UpdateProduct(pid int, stock int) (models.ProductResponse, error) {

	result, err := i.repository.CheckProduct(pid)
	if err != nil {
		return models.ProductResponse{}, err
	}

	if !result {
		return models.ProductResponse{}, errors.New("there is no inventory as you mentioned")
	}

	newcat, err := i.repository.UpdateProduct(pid, stock)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return newcat, err
}
