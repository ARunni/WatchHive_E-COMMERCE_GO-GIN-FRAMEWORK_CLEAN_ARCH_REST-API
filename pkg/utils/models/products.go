package models

type ProductResponse struct {
	ID          int    `json:"id"`
	CategoryID  int    `json:"category_id"`
	ProductName string `json:"productname"`
	Color       string `json:"color"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
}

type ProductUpdate struct {
	Productid int `json:"product_id"`
	Stock     int `json:"stock"`
}

type Product struct {
	ID          uint   `json:"id"`
	CategoryID  int    `json:"category_id"`
	ProductName string `json:"productname"`
	Color       string `json:"color"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
}

type AddProducts struct {
	ID          uint    `json:"id"`
	CategoryID  int     `json:"category_id"`
	ProductName string  `json:"product_name"`
	Color       string  `json:"color"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}

type EditProductDetials struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID int     `json:"category_id"`
	Color      string  `json:"color"`
}

type ProductUserResponse struct {
	ID          uint   `json:"id"`
	CategoryID  int    `json:"category_id"`
	Category    string `json:"category"`
	ProductName string `json:"product_name"`
	Color       string `json:"color"`
	Price       int    `json:"price"`
	Url         string `json:"image"`
	Stock       int    `json:"stock"`
}
type ProductEdit struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	CategoryID  uint    `json:"category_id"`
	ProductName string  `json:"product_name"`
	Color       string  `json:"color" gorm:"color:6;default:'Black';Check:color IN ('Black', 'Blue', 'Red', 'Green','RoseGold');"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}
