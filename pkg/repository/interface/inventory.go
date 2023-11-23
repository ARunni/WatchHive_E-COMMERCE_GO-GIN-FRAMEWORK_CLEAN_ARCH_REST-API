package interfaces

import (
	"WatchHive/pkg/domain"
	"WatchHive/pkg/utils/models"
)

type InventoryRepository interface{
	
	AddInventory(inventory models.AddInventories) (models.InventoryResponse, error)
	ListProducts(int, int) ([]models.InventoryUserResponse, error)
	EditInventory(domain.Inventory, int) (domain.Inventory, error)
	DeleteInventory(id string) error
	CheckInventory(pid int) (bool, error)
	UpdateInventory(pid int, stock int) (models.InventoryResponse, error)
	
}