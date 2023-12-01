package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(Db *gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		DB: Db,
	}
}


func (or *orderRepository) GetAllPaymentOption() ([]models.PaymentDetails, error) {
	var paymentMethods []models.PaymentDetails
	err := or.DB.Raw("SELECT * FROM payment_methods").Scan(&paymentMethods).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}

	return paymentMethods, nil

}
func (or *orderRepository) GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	var addressId int
	if err := or.DB.Raw("SELECT address_id FROM orders WHERE id =?", orderID).Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("first in orders")
	}
	if err := or.DB.Raw("SELECT * FROM addresses WHERE id=?", addressId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("second  in address")
	}
	return addressInfoResponse, nil
}
func (or *orderRepository) GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := or.DB.Raw("SELECT id,final_price,shipment_status,payment_status FROM orders WHERE id = ?", orderID).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil
}

func (or *orderRepository) GetProductsInCart(cart_id int) ([]int, error) {

	var cart_products []int

	if err := or.DB.Raw("select product_id from cart_items where cart_id=?", cart_id).Scan(&cart_products).Error; err != nil {
		return []int{}, err
	}

	return cart_products, nil

}
func (or *orderRepository) FindProductNames(product_id int) (string, error) {

	var product_name string

	if err := or.DB.Raw("select name from products where id=?", product_id).Scan(&product_name).Error; err != nil {
		return "", err
	}

	return product_name, nil

}

func (or *orderRepository) FindCartQuantity(cart_id, product_id int) (int, error) {

	var quantity int

	if err := or.DB.Raw("select quantity from cart_items where cart_id=$1 and product_id=$2", cart_id, product_id).Scan(&quantity).Error; err != nil {
		return 0, err
	}

	return quantity, nil

}

func (or *orderRepository) FindPrice(product_id int) (float64, error) {

	var price float64

	if err := or.DB.Raw("select price from products where id=?", product_id).Scan(&price).Error; err != nil {
		return 0, err
	}

	return price, nil

}
func (or *orderRepository) FindStock(id int) (int, error) {
	var stock int
	err := or.DB.Raw("SELECT stock FROM prodcuts WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}
