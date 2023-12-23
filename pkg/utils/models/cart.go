package models

type CartResponse struct {
	UserName   string
	TotalPrice float64
	Cart       []Cart
}
type CartTotal struct {
	UserName       string  `json:"user_name"`
	TotalPrice     float64 `json:"total_price"`
	FinalPrice     float64 `json:"final_price"`
	DiscountReason string
}
type Cart struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type Carts struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
}
type GetCartResponse struct {
	ID   int
	Data []GetCart
}
type GetCart struct {
	ID          int     `json:"id" gorm:"primarykey;not null"`
	ProductId   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Product     Product `json:"-" gorm:"foreignkey:ProductId"`
	Quantity    float64 `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type AddCart struct {
	UserID    int     `json:"id"`
	ProductID int     `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

type RemoveFromCart struct {
	UserID    int `json:"id"`
	ProductID int `json:"product_id"`
}
type RemoveFromCartR struct {
	ProductID int `json:"product_id"`
}

type AddCartR struct {
    ProductID int     `json:"product_id"`
    Quantity  float64 `json:"quantity"`
}