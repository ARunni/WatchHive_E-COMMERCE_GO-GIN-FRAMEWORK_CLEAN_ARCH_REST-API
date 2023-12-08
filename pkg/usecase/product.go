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
	cat        rep.CategoryRepository
}

func NewProductUseCase(repo rep.ProductRepository, h helper.Helper, cat rep.CategoryRepository) interfaces.ProductUseCase {
	return &productUseCase{
		repository: repo,
		helper:     h,
		cat:        cat,
	}

}

func (i *productUseCase) AddProduct(product models.AddProducts, file *multipart.FileHeader) (models.ProductResponse, error) {
	if product.ProductName == "" {
		return models.ProductResponse{}, errors.New("product name cannot be empty")

	}
	if product.Color == "" {
		return models.ProductResponse{}, errors.New("color cannot be empty")
	}
	if product.CategoryID <= 0 || product.Price <= 0 || product.Stock <= 0 {
		err := errors.New("enter valid values")
		return models.ProductResponse{}, err
	}
	ok, err := i.cat.CheckCategory(product.CategoryID)
	if err != nil {
		return models.ProductResponse{}, err
	}
	if !ok {
		return models.ProductResponse{}, errors.New("category id not available")
	}
	ok = i.repository.CheckProductAndCat(product.ProductName, product.CategoryID)

	if ok {
		return models.ProductResponse{}, errors.New("already exist")
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
	if pageNo <= 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageList
	productList, err := i.repository.ListProducts(pageList, offset)
	if err != nil {
		return []models.ProductUserResponse{}, err
	}
	return productList, nil
}

func (usecase *productUseCase) EditProduct(product domain.Product, id int) (domain.Product, error) {

	if product.ID <= 0 || product.CategoryID <= 0 || product.Price <= 0 || product.Stock <= 0 {
		err := errors.New("enter valid values")
		return domain.Product{}, err
	}
	if product.ProductName == "" {
		return domain.Product{}, errors.New("product name cannot be empty")
	}
	if product.Color == "" {
		return domain.Product{}, errors.New("color cannot be empty")
	}
	modProduct, err := usecase.repository.EditProduct(product, id)
	if err != nil {
		return domain.Product{}, err
	}
	return modProduct, nil
}

func (pu *productUseCase) DeleteProduct(productID string) error {

	pID, cErr := strconv.Atoi(productID)

	if cErr != nil || pID <= 0 {
		return errors.New("invalid data")

	}
	ok, err := pu.repository.CheckProduct(pID)
	if err != nil {
		return errors.New("error in searching")
	}
	if !ok {
		return errors.New("product does not exist")
	}
	err = pu.repository.DeleteProduct(productID)
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
		return models.ProductResponse{}, errors.New("error occured during the search")
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
