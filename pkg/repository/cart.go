package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(DB *gorm.DB) interfaces.CartRepository {
	return &cartRepository{
		DB: DB,
	}
}

func (cr *cartRepository) AddToCart(userId int, productId int, Quantity int, productprice float64) error {

	query := "INSERT INTO carts (user_id,product_id,quantity,total_price) VALUES (?,?,?,?)"
	if err := cr.DB.Exec(query, userId, productId, Quantity, productprice).Error; err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) CheckProductAvailable(product_id int) (bool, error) {
	var count int
	querry := "SELECT COUNT(*) FROM products where id = ?"

	err := cr.DB.Raw(querry, product_id).Scan(&count).Error
	if err != nil {
		return false, errors.New("product does not exist")
	}
	if count < 1 {
		return false, errors.New("product does not exist")
	}
	return true, nil
}

func (cr *cartRepository) CheckStock(product_id int) (int, error) {
	qurry := "SELECT stock FROM products where id = ?"
	var stock int
	err := cr.DB.Raw(qurry, product_id).Scan(&stock).Error
	if err != nil {
		return 0, errors.New("error in getting stock")
	}
	return stock, nil
}

func (cr *cartRepository) QuantityOfProductInCart(userId int, productId int) (int, error) {
	var productQty int
	querry := "SELECT quantity FROM carts WHERE user_id = ? AND product_id = ?"
	err := cr.DB.Raw(querry, userId, productId).Scan(&productQty).Error
	if err != nil {
		return 0, errors.New("error in getting quantity")
	}
	return productQty, nil
}

func (cr *cartRepository) GetTotalPriceFromCart(userID int) (float64, error) {
	var totalPrice float64
	err := cr.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&totalPrice).Error
	if err != nil {
		return 0.0, err
	}
	return totalPrice, nil

}
func (cr *cartRepository) UpdateCart(quantity int, price float64, userID int, product_id int) error {
	if err := cr.DB.Exec(`UPDATE carts
	SET quantity = ?, total_price = ? 
	WHERE user_id = ? and product_id = ?`, quantity, price, product_id, userID).Error; err != nil {
		return err
	}

	return nil

}

func (cr *cartRepository) TotalPriceForProductInCart(userID int, productID int) (float64, error) {

	var totalPrice float64
	if err := cr.DB.Raw("SELECT SUM(total_price) as total_price FROM carts  WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}

func (cr *cartRepository) DisplayCart(userID int) ([]models.Cart, error) {

	var count int
	if err := cr.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}

	var cartResponse []models.Cart
	qurry := `SELECT carts.user_id,users.name as name,carts.product_id,
				products.product_name as product_name,carts.quantity,carts.total_price 
				from carts INNER JOIN users ON carts.user_id = users.id 
				INNER JOIN products ON carts.product_id = products.id WHERE user_id = ?`

	if err := cr.DB.Raw(qurry, userID).First(&cartResponse).Error; err != nil {
		return []models.Cart{}, err
	}
	return cartResponse, nil

}

func (cr *cartRepository) GetTotalPrice(userID int) (models.CartTotal, error) {

	var cartTotal models.CartTotal
	err := cr.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.TotalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = cr.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.FinalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = cr.DB.Raw("SELECT name as user_name FROM users WHERE id = ?", userID).Scan(&cartTotal.UserName).Error
	if err != nil {
		return models.CartTotal{}, err
	}

	return cartTotal, nil

}

//jhhjgfhajghkkhjg

func (cr *cartRepository) UpdateProductQuantityCart(cart models.AddCart) error {

	querry := `	UPDATE carts
	SET quantity = $1, total_price = $1 * (select price from products where id = $3)
	WHERE user_id= $2 AND product_id= $3`

	err := cr.DB.Exec(querry, cart.Quantity, cart.UserID, cart.ProductID).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) RemoveFromCart(cart models.RemoveFromCart) error {
	querry := `	DELETE  FROM carts WHERE product_id = ? AND user_id = ?`
	err := cr.DB.Exec(querry, cart.ProductID, cart.UserID).Error
	if err != nil {
		return errors.New("error at database")
	}
	return nil
}

func (cr *cartRepository) CheckCart(userID int) (bool, error) {
	var count int
	querry := `	SELECT COUNT(*) FROM carts WHERE user_id = ?`
	err := cr.DB.Raw(querry, userID).Scan(&count).Error
	if err != nil {
		return false, errors.New("no cart found")
	}
	if count < 0 {
		return false, errors.New("no cart found")
	}
	return true, nil
}
