package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	DB *gorm.DB
}

func NewInventoryRepository(DB *gorm.DB) interfaces.InventoryRepository {
	return &inventoryRepository{
		DB: DB,
	}
}

func (i *inventoryRepository) AddInventory(inventory models.AddInventories) (models.InventoryResponse, error) {

	var count int64
	i.DB.Model(&models.Inventory{}).Where("product_name = ? AND category_id = ?", inventory.ProductName, inventory.CategoryID).Count(&count)
	if count > 0 {

		return models.InventoryResponse{}, errors.New("product already exists in the database")
	}

	if inventory.Stock < 0 || inventory.Price < 0 {
		return models.InventoryResponse{}, errors.New("stock and price cannot be negative")
	}

	query := `
        INSERT INTO inventories (category_id, product_name, color, stock, price)
        VALUES (?, ?, ?, ?, ?);
    `
	err := i.DB.Exec(query, inventory.CategoryID, inventory.ProductName, inventory.Color, inventory.Stock, inventory.Price).Error
	if err != nil {
		return models.InventoryResponse{}, err
	}

	var inventoryResponse models.InventoryResponse

	return inventoryResponse, nil
}

func (prod *inventoryRepository) ListProducts(pageList, offset int) ([]models.InventoryUserResponse, error) {

	var product_list []models.InventoryUserResponse

	query := "SELECT i.id,i.category_id,c.category,i.product_name,i.color,i.price FROM inventories i INNER JOIN categories c ON i.category_id = c.id LIMIT $1 OFFSET $2"
	err := prod.DB.Raw(query, pageList, offset).Scan(&product_list).Error

	if err != nil {
		return []models.InventoryUserResponse{}, errors.New("error checking Product details")
	}

	return product_list, nil
}

func (db *inventoryRepository) EditInventory(inventory domain.Inventory, id int) (domain.Inventory, error) {

	var modInventory domain.Inventory

	query := "UPDATE inventories SET category_id = ?, product_name = ?, color = ?, stock = ?, price = ? WHERE id = ?"

	if err := db.DB.Exec(query, inventory.CategoryID, inventory.ProductName, inventory.Color, inventory.Stock, inventory.Price, id).Error; err != nil {
		return domain.Inventory{}, err
	}

	if err := db.DB.First(&modInventory, id).Error; err != nil {
		return domain.Inventory{}, err
	}
	return modInventory, nil
}
func (i *inventoryRepository) DeleteInventory(inventoryID string) error {

	id, err := strconv.Atoi(inventoryID)
	if err != nil {
		return errors.New("converting into integet is not happened")
	}

	result := i.DB.Exec("DELETE FROM inventories WHERE id = ?", id)

	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

func (i *inventoryRepository) CheckInventory(pid int) (bool, error) {
	var k int
	err := i.DB.Raw("SELECT COUNT(*) FROM inventories WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return false, err
	}

	if k == 0 {
		return false, err
	}

	return true, err
}

func (i *inventoryRepository) UpdateInventory(pid int, stock int) (models.InventoryResponse, error) {

	// Check the database connection
	if i.DB == nil {
		return models.InventoryResponse{}, errors.New("database connection is nil")
	}

	// Update the
	if err := i.DB.Exec("UPDATE inventories SET stock = stock + $1 WHERE id= $2", stock, pid).Error; err != nil {
		return models.InventoryResponse{}, err
	}

	// Retrieve the update
	var newdetails models.InventoryResponse
	var newstock int
	if err := i.DB.Raw("SELECT stock FROM inventories WHERE id=?", pid).Scan(&newstock).Error; err != nil {
		return models.InventoryResponse{}, err
	}
	newdetails.ProductID = pid
	newdetails.Stock = newstock

	return newdetails, nil
}
