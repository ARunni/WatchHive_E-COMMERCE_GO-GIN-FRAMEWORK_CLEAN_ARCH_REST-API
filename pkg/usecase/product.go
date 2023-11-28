package usecase

import (
	"WatchHive/pkg/domain"
	helper "WatchHive/pkg/helper/interface"
	rep "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"mime/multipart"
	"strconv"
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

func (i *productUseCase) AddProduct(product models.AddProducts, file *multipart.FileHeader) (models.ProductResponse, error) {

	if product.CategoryID < 0 || product.Price < 0 || product.Stock < 0 {
		err := errors.New("enter valid values")
		return models.ProductResponse{}, err
	}

	url, err := i.helper.AddImageToAwsS3(file)
	if err != nil {
		return models.ProductResponse{}, err
	}

	ProductResponse, err := i.repository.AddProduct(product, url)
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

	if product.ID == 0 || product.CategoryID == 0 || product.Price < 0 || product.Stock < 0 {
		err := errors.New("enter valid values")
		return domain.Product{}, err
	}

	modProduct, err := usecase.repository.EditProduct(product, id)
	if err != nil {
		return domain.Product{}, err
	}
	return modProduct, nil
}

func (usecase *productUseCase) DeleteProduct(productID string) error {

	id, cErr := strconv.Atoi(productID)

	if cErr != nil || id <= 0 {
		return cErr

	}

	err := usecase.repository.DeleteProduct(productID)
	if err != nil {
		return err
	}
	return nil
}

func (i productUseCase) UpdateProduct(pid int, stock int) (models.ProductResponse, error) {

	if pid <= 0 || stock <= 0 {
		return models.ProductResponse{}, errors.New("invalid product id or stock")
	}

	result, err := i.repository.CheckProduct(pid)
	if err != nil {
		return models.ProductResponse{}, err
	}

	if !result {
		return models.ProductResponse{}, errors.New("there is no product as you mentioned")
	}

	newcat, err := i.repository.UpdateProduct(pid, stock)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return newcat, err
}
