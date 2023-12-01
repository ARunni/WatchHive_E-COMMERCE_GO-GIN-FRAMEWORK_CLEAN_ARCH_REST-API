package domain

type Product struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName string   `json:"product_name"`
	Color       string   `json:"color" gorm:"color:6;default:'Black';Check:color IN ('Black', 'Blue', 'Red', 'Green','RoseGold');"`
	Stock       int      `json:"stock"`
	Price       float64  `json:"price"`
}

type ProductImage struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	ProductId int  `json:"product_id"`
	// Product   Product `json:"product" gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Url string `json:"url"`
}
