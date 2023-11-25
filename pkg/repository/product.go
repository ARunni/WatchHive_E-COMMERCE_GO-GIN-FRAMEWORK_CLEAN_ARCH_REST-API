package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (i *productRepository) AddProduct(product models.AddProducts) (models.ProductResponse, error) {

	var count int64
	i.DB.Model(&models.Product{}).Where("product_name = ? AND category_id = ?", product.ProductName, product.CategoryID).Count(&count)
	if count > 0 {

		return models.ProductResponse{}, errors.New("product already exists in the database")
	}

	if product.Stock < 0 || product.Price < 0 {
		return models.ProductResponse{}, errors.New("stock and price cannot be negative")
	}

	query := `
        INSERT INTO products (category_id, product_name, color, stock, price)
        VALUES (?, ?, ?, ?, ?);
    `
	err := i.DB.Exec(query, product.CategoryID, product.ProductName, product.Color, product.Stock, product.Price).Error
	if err != nil {
		return models.ProductResponse{}, err
	}
	// getting inserted product detailsS
	var productResponse models.ProductResponse

	query = "SELECT id,category_id,product_name,color,stock,price FROM products where  product_name = ? AND category_id = ?"
	errr := i.DB.Raw(query, product.ProductName, product.CategoryID).Scan(&productResponse).Error

	if errr != nil {
		return productResponse, errors.New("error checking Product details")
	}

	return productResponse, nil
}

func (prod *productRepository) ListProducts(pageList, offset int) ([]models.ProductUserResponse, error) {

	var product_list []models.ProductUserResponse

	query := "SELECT i.id,i.category_id,c.category,i.product_name,i.color,i.price FROM products i INNER JOIN categories c ON i.category_id = c.id LIMIT $1 OFFSET $2"
	err := prod.DB.Raw(query, pageList, offset).Scan(&product_list).Error

	if err != nil {
		return []models.ProductUserResponse{}, errors.New("error checking Product details")
	}

	return product_list, nil
}

func (db *productRepository) EditProduct(product domain.Product, id int) (domain.Product, error) {

	var modProduct domain.Product

	query := "UPDATE products SET category_id = ?, product_name = ?, color = ?, stock = ?, price = ? WHERE id = ?"

	if err := db.DB.Exec(query, product.CategoryID, product.ProductName, product.Color, product.Stock, product.Price, id).Error; err != nil {
		return domain.Product{}, err
	}

	if err := db.DB.First(&modProduct, id).Error; err != nil {
		return domain.Product{}, err
	}
	return modProduct, nil
}
func (i *productRepository) DeleteProduct(productID string) error {

	id, err := strconv.Atoi(productID)
	if err != nil {
		return errors.New("converting into integet is not happened")
	}

	result := i.DB.Exec("DELETE FROM products WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

func (i *productRepository) CheckProduct(pid int) (bool, error) {
	var k int
	err := i.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return false, err
	}

	if k == 0 {
		return false, err
	}

	return true, err
}

func (i *productRepository) UpdateProduct(pid int, stock int) (models.ProductResponse, error) {

	// Check the database connection
	if i.DB == nil {
		return models.ProductResponse{}, errors.New("database connection is nil")
	}

	// Update the
	if err := i.DB.Exec("UPDATE products SET stock = stock + $1 WHERE id= $2", stock, pid).Error; err != nil {
		return models.ProductResponse{}, err
	}

	// Retrieve the update
	var newdetails models.ProductResponse
	var newstock int
	if err := i.DB.Raw("SELECT stock FROM products WHERE id=?", pid).Scan(&newstock).Error; err != nil {
		return models.ProductResponse{}, err
	}
	newdetails.ProductID = pid
	newdetails.Stock = newstock

	return newdetails, nil
}
