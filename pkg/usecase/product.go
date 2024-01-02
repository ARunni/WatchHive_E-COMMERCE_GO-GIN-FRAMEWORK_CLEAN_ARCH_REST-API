package usecase

import (
	helper "WatchHive/pkg/helper/interface"
	rep "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
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
		return models.ProductResponse{}, errors.New("product name " + errmsg.ErrFieldEmpty)

	}
	if product.Color == "" {
		return models.ProductResponse{}, errors.New("color " + errmsg.ErrFieldEmpty)
	}
	if product.CategoryID <= 0 || product.Price <= 0 || product.Stock <= 0 {
		err := errors.New(errmsg.ErrInvalidData)
		return models.ProductResponse{}, err
	}
	ok, err := i.cat.CheckCategory(product.CategoryID)
	if err != nil {
		return models.ProductResponse{}, err
	}
	if !ok {
		return models.ProductResponse{}, errors.New(errmsg.ErrInvalidCId)
	}
	ok = i.repository.CheckProductAndCat(product.ProductName, product.CategoryID)

	if ok {
		return models.ProductResponse{}, errors.New(errmsg.ErrExistTrue)
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

func (usecase *productUseCase) EditProduct(product models.ProductEdit) (models.ProductUserResponse, error) {

	if product.ID <= 0 || product.CategoryID <= 0 || product.Price <= 0 || product.Stock <= 0 {
		err := errors.New(errmsg.ErrInvalidData)
		return models.ProductUserResponse{}, err
	}
	if product.ProductName == "" {
		return models.ProductUserResponse{}, errors.New("product name "+ errmsg.ErrFieldEmpty)
	}
	if product.Color == "" {
		return models.ProductUserResponse{}, errors.New("color " +errmsg.ErrFieldEmpty )
	}
	modProduct, err := usecase.repository.EditProduct(product)
	if err != nil {
		return models.ProductUserResponse{}, err
	}
	return modProduct, nil
}

func (pu *productUseCase) DeleteProduct(productID string) error {

	pID, cErr := strconv.Atoi(productID)

	if cErr != nil || pID <= 0 {
		return errors.New(errmsg.ErrInvalidData)

	}
	ok, err := pu.repository.CheckProduct(pID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("product "+ errmsg.ErrNotExist)
	}
	err = pu.repository.DeleteProduct(productID)
	if err != nil {
		return err
	}
	return nil
}

func (i productUseCase) UpdateProduct(pid int, stock int) (models.ProductResponse, error) {

	if pid <= 0 || stock <= 0 {
		return models.ProductResponse{}, errors.New(errmsg.ErrInvalidData)
	}

	result, err := i.repository.CheckProduct(pid)
	if err != nil {
		return models.ProductResponse{}, err
	}

	if !result {
		return models.ProductResponse{}, errors.New("product " + errmsg.ErrNotExist)
	}

	newcat, err := i.repository.UpdateProduct(pid, stock)
	if err != nil {
		return models.ProductResponse{}, err
	}

	return newcat, err
}
