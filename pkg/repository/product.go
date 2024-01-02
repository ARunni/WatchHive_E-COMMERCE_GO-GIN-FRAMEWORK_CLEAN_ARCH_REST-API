package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
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

func (i *productRepository) AddProduct(product models.AddProducts, url string) (models.ProductResponse, error) {

	var count int64
	i.DB.Model(&models.Product{}).Where("product_name = ? AND category_id = ?", product.ProductName, product.CategoryID).Count(&count)
	if count > 0 {

		return models.ProductResponse{}, errors.New(errmsg.ErrProductExistTrue)
	}

	if product.Stock < 0 || product.Price < 0 {
		return models.ProductResponse{}, errors.New("stock and price" + errmsg.ErrDataNegative)
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
		return productResponse, errors.New(errmsg.ErrGetDB)
	}
	//Adding url to image table

	queryimage := "INSERT INTO product_images (product_id, url) VALUES (?, ?)"

	imgErr := i.DB.Exec(queryimage, productResponse.ID, url).Error
	if err != nil {

		return models.ProductResponse{}, imgErr
	}
	return productResponse, nil
}

func (prod *productRepository) ListProducts(pageList, offset int) ([]models.ProductUserResponse, error) {

	var product_list []models.ProductUserResponse

	query := `SELECT p.id, p.category_id, c.category, p.product_name, p.color, p.price, i.url, p.stock
	FROM products p
	INNER JOIN categories c ON p.category_id = c.id
	LEFT JOIN product_images i ON p.id = i.product_id
	LIMIT $1 OFFSET $2`

	err := prod.DB.Raw(query, pageList, offset).Scan(&product_list).Error

	if err != nil {
		return []models.ProductUserResponse{}, errors.New(errmsg.ErrGetDB)
	}

	return product_list, nil
}

func (db *productRepository) EditProduct(product models.ProductEdit) (models.ProductUserResponse, error) {

	var modProduct models.ProductUserResponse

	query := "UPDATE products SET category_id = ?, product_name = ?, color = ?, stock = ?, price = ? WHERE id = ?"

	if err := db.DB.Exec(query, product.CategoryID, product.ProductName, product.Color, product.Stock, product.Price, product.ID).Error; err != nil {
		return models.ProductUserResponse{}, err
	}

	if err := db.DB.Raw("select * from products where id = ?", product.ID).Scan(&modProduct).Error; err != nil {
		return models.ProductUserResponse{}, err
	}
	return modProduct, nil
}
func (i *productRepository) DeleteProduct(productID string) error {

	id, err := strconv.Atoi(productID)
	if err != nil {

		return errors.New(errmsg.ErrDatatypeConversion)
	}

	result := i.DB.Exec("DELETE FROM products WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New(errmsg.ErrUserExistFalse)
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

func (i *productRepository) CheckProductAndCat(prdt string, cat int) bool {
	var count int
	err := i.DB.Raw("SELECT COUNT(*) FROM products WHERE category_id=? And product_name=?", cat, prdt).Scan(&count).Error
	if err != nil {
		return false
	}
	if count <= 0 {
		return false
	}
	return true
}

func (i *productRepository) UpdateProduct(pid int, stock int) (models.ProductResponse, error) {

	// Check the database connection
	if i.DB == nil {
		return models.ProductResponse{}, errors.New(errmsg.ErrDbConnect)
	}

	// Update the
	if err := i.DB.Exec("UPDATE products SET stock = stock + $1 WHERE id= $2", stock, pid).Error; err != nil {
		return models.ProductResponse{}, err
	}

	// Retrieve the update
	var newdetails models.ProductResponse
	// var newstock int
	if err := i.DB.Raw("SELECT * FROM products WHERE id=?", pid).Scan(&newdetails).Error; err != nil {
		return models.ProductResponse{}, err
	}
	// newdetails.ID = pid
	// newdetails.Stock = newstock

	return newdetails, nil
}

func (cr *productRepository) CheckProductAvailable(product_id int) (bool, error) {
	var count int
	querry := "SELECT COUNT(*) FROM products where id = ?"

	err := cr.DB.Raw(querry, product_id).Scan(&count).Error
	if err != nil {
		return false, errors.New(errmsg.ErrProductExist)
	}
	if count < 1 {
		return false, errors.New(errmsg.ErrProductExist)
	}
	return true, nil
}
func (cr *productRepository) GetPriceOfProduct(product_id int) (float64, error) {
	qurry := "SELECT price FROM products where id = ?"
	var price float64
	err := cr.DB.Raw(qurry, product_id).Scan(&price).Error

	if err != nil {
		return 0, errors.New(errmsg.ErrGetDB)
	}
	return price, nil
}
