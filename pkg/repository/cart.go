package repository

// import (
// 	interfaces "WatchHive/pkg/repository/interface"
// 	"WatchHive/pkg/utils/models"

// 	"gorm.io/gorm"
// )

// type cartRepository struct {
// 	DB *gorm.DB
// }

// func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
// 	return &cartRepository{
// 		DB: DB,
// 	}
// }

// func (cr *cartRepository) AddToCart(userId int, productId int, Quantity int, productprice float64) error {

// 	query := "INSERT INTO carts (user_id,product_id,quantity,total_price) VALUES (?,?,?,?)"
// 	if err := cr.DB.Exec(query, userId, productId, Quantity, productprice).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (cr *cartRepository) DisplayCart(userID int) ([]models.Cart, error) {
// 	var count int
// 	qry := "SELECT COUNT(*) FROM carts WHERE user_id = ? "
// 	if err := cr.DB.Raw(qry, userID).Scan(&count).Error; err != nil {
// 		return []models.Cart{}, err
// 	}
// 	qry = "SELECT carts.user_id,users.name as user_name,carts.product_id,products.product_name as product_name,carts.quantity,carts.total_price from carts INNER JOIN users ON carts.user_id = users.id INNER JOIN products ON carts.product_id = products.id WHERE user_id = ?"
// 	var cartResponse []models.Cart

// 	if err := cr.DB.Raw(qry, userID).First(&cartResponse).Error; err != nil {
// 		return []models.Cart{}, err
// 	}
// 	return cartResponse, nil
// }
