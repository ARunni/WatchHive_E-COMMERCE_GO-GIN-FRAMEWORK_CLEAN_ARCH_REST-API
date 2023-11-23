package models

type InventoryResponse struct {
	ProductID   int    `json:"id"`
	CategoryID  int    `json:"category_id"`
	ProductName string `json:"productname"`
	Color       string `json:"color"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
}

type InventoryUpdate struct {
	Productid int `json:"product_id"`
	Stock     int `json:"stock"`
}

type Inventory struct {
	ID          uint   `json:"id"`
	CategoryID  int    `json:"category_id"`
	ProductName string `json:"productname"`
	Color       string `json:"color"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
}

type AddInventories struct {
	ID          uint    `json:"id"`
	CategoryID  int     `json:"category_id"`
	ProductName string  `json:"product_name"`
	Color       string  `json:"color"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}

type EditInventoryDetials struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID int     `json:"category_id"`
	Color      string  `json:"color"`
}

type InventoryUserResponse struct {
	ID          uint   `json:"id"`
	CategoryID  int    `json:"category_id"`
	Category    string `json:"category"`
	ProductName string `json:"productname"`
	Color       string `json:"color"`
	Price       int    `json:"price"`
}
