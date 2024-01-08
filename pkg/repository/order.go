package repository

import (
	"WatchHive/pkg/domain"
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/errmsg"
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
		return models.AddressInfoResponse{}, errors.New(errmsg.ErrGetDB)
	}
	if err := or.DB.Raw("SELECT * FROM addresses WHERE id=?", addressId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New(errmsg.ErrGetDB)
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

func (or *orderRepository) CheckOrderID(orderId int) (bool, error) {
	var count int
	err := or.DB.Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderId).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (or *orderRepository) OrderExist(orderID int) (bool, error) {
	var count int
	err := or.DB.Raw("SELECT count(*) FROM orders WHERE id = ?", orderID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	if count <= 0 {
		return false, nil
	}
	return true, nil
}

func (or *orderRepository) GetShipmentStatus(orderID int) (string, error) {
	var status string
	err := or.DB.Raw("SELECT shipment_status FROM orders WHERE id= ?", orderID).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

// pay
func (or *orderRepository) GetPaymentType(orderID int) (int, error) {
	var status int
	err := or.DB.Raw("SELECT payment_method_id FROM orders WHERE id= ?", orderID).Scan(&status).Error
	if err != nil {
		return 0, err
	}
	return status, nil
}

// pay

func (or *orderRepository) UpdateOrder(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET Shipment_status = 'processing' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) AddOrderProducts(order_id int, cart []models.Cart) error {
	query := `
    INSERT INTO order_items (order_id,product_id,quantity,total_price)
    VALUES (?, ?, ?, ?) `
	for _, v := range cart {
		var productID int
		if err := or.DB.Raw("SELECT id FROM products WHERE product_name = $1", v.ProductName).Scan(&productID).Error; err != nil {
			return err
		}
		if err := or.DB.Exec(query, order_id, productID, v.Quantity, v.TotalPrice).Error; err != nil {
			return err
		}
	}
	return nil
}

func (or *orderRepository) GetBriefOrderDetails(orderID int) (models.OrderSuccessResponse, error) {
	var orderSuccessResponse models.OrderSuccessResponse
	err := or.DB.Raw(`SELECT id as order_id,shipment_status FROM orders WHERE id = ?`, orderID).Scan(&orderSuccessResponse).Error
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	return orderSuccessResponse, nil
}

func (or *orderRepository) OrderItems(ob models.OrderIncoming, price float64) (int, error) {
	var id int
	query := `
    INSERT INTO orders (created_at , user_id , address_id , payment_method_id , final_price)
    VALUES (NOW(),?, ?, ?, ?)
    RETURNING id`
	or.DB.Raw(query, ob.UserID, ob.AddressID, ob.PaymentID, price).Scan(&id)
	return id, nil
}

func (or *orderRepository) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDetails []models.OrderDetails
	err := or.DB.Raw("SELECT id as order_id,final_price,shipment_status,payment_status FROM orders WHERE user_id = ? order by id desc LIMIT ? OFFSET ? ", userId, count, offset).Scan(&orderDetails).Error

	if err != nil {
		return []models.FullOrderDetails{}, err
	}

	var fullOrderDetails []models.FullOrderDetails
	for _, od := range orderDetails {
		var orderProductDetails []models.OrderProductDetails
		err := or.DB.Raw(`SELECT
		order_items.product_id,
		products.product_name AS product_name,
		order_items.quantity,
		order_items.total_price
	    FROM
		order_items
	    INNER JOIN
		products ON order_items.product_id = products.id
	    WHERE
		order_items.order_id = $1 `, od.OrderId).Scan(&orderProductDetails).Error
		if err != nil {
			return []models.FullOrderDetails{}, err
		}
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{OrderDetails: od, OrderProductDetails: orderProductDetails})
	}
	return fullOrderDetails, nil
}

func (or *orderRepository) UserOrderRelationship(orderID int, userID int) (int, error) {

	var testUserID int
	err := or.DB.Raw("select user_id from orders where id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}
func (or *orderRepository) GetPaymentStatus(orderID int) (string, error) {
	var paymentStatus string
	if err := or.DB.Raw("select payment_status from orders where id = ?", orderID).Scan(&paymentStatus).Error; err != nil {
		return "", errors.New(errmsg.ErrGetDB)
	}
	return paymentStatus, nil
}
func (or *orderRepository) GetProductDetailsFromOrders(orderID int) ([]models.OrderProducts, error) {
	var OrderProductDetails []models.OrderProducts
	if err := or.DB.Raw("SELECT product_id,quantity as stock FROM order_items WHERE order_id = ?", orderID).Scan(&OrderProductDetails).Error; err != nil {
		return []models.OrderProducts{}, err
	}
	return OrderProductDetails, nil
}
func (or *orderRepository) CancelOrders(orderID int) error {
	status := "cancelled"
	err := or.DB.Exec("UPDATE orders SET shipment_status = ? , approval='false' WHERE id = ? ", status, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = or.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ? ", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 3 || paymentMethod == 2 {
		err = or.DB.Exec("UPDATE orders SET payment_status = 'credited to wallet' WHERE id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	if paymentMethod == 1 {
		err = or.DB.Exec("UPDATE orders SET payment_status = 'cod canceled' WHERE id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (or *orderRepository) ReturnOrderCod(orderId int) error {

	shipStatus := "returned"
	payStatus := "processing"
	err := or.DB.Exec("UPDATE orders SET shipment_status = ? , approval='false',payment_status = ? WHERE id = ? ", shipStatus, payStatus, orderId).Error
	if err != nil {
		return err
	}
	return nil

}
func (or *orderRepository) ReturnOrderRazorPay(orderId int) error {

	shipStatus := "returned"
	payStatus := "credited to wallet"
	err := or.DB.Exec("UPDATE orders SET shipment_status = ? , approval='false',payment_status = ? WHERE id = ? ", shipStatus, payStatus, orderId).Error
	if err != nil {
		return err
	}
	return nil

}

func (or *orderRepository) UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {

	for _, od := range orderProducts {

		var quantity int
		if err := or.DB.Raw("SELECT stock FROM products WHERE id = ?", od.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}

		od.Stock += quantity
		if err := or.DB.Exec("UPDATE products SET stock = ? WHERE id = ?", od.Stock, od.ProductId).Error; err != nil {
			return err
		}
	}
	return nil

}
func (or *orderRepository) GetAllOrdersAdmin(offset, count int) ([]models.CombinedOrderDetails, error) {

	var orderDatails []models.CombinedOrderDetails
	querry := `
SELECT orders.id as order_id,orders.final_price,
orders.shipment_status,orders.payment_status,
users.name,users.email,users.phone,
addresses.house_name,addresses.street,
addresses.city,addresses.state,
addresses.pin 
FROM orders INNER JOIN users 
ON orders.user_id = users.id INNER JOIN addresses 
ON orders.address_id = addresses.id order by orders.id desc limit ? offset ?`

	err := or.DB.Raw(querry, count, offset).Scan(&orderDatails).Error
	if err != nil {
		return []models.CombinedOrderDetails{}, nil
	}
	return orderDatails, nil

}

func (or *orderRepository) ApproveOrder(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET shipment_status = 'shipped' , approval = 'true' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) ApproveCodPaid(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET shipment_status = 'delivered' , approval = 'true', payment_status = 'PAID' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func (or *orderRepository) ApproveRazorPaid(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET shipment_status = 'shipped' , approval = 'true', payment_status = 'PAID' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func (or *orderRepository) ApproveRazorDelivered(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET shipment_status = 'delivered' , approval = 'true', payment_status = 'PAID' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) ApproveCodReturn(orderID int) error {
	err := or.DB.Exec("UPDATE orders SET   approval = 'true', payment_status = 'added_to_wallet' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func (or *orderRepository) UpdateStockOfProduct(orderProducts []models.OrderProducts) error {
	for _, ok := range orderProducts {
		var quantity int
		if err := or.DB.Raw("SELECT stock FROM products WHERE id = ?", ok.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}
		ok.Stock += quantity
		if err := or.DB.Exec("UPDATE products SET stock  = ? WHERE id = ?", ok.Stock, ok.ProductId).Error; err != nil {
			return err
		}
	}
	return nil
}

//razor

func (repo *orderRepository) GetOrder(orderId int) (domain.Order, error) {
	var body domain.Order
	query := `
		select * from orders
		where id = $1
	`
	if err := repo.DB.Raw(query, orderId).Scan(&body).Error; err != nil {
		return domain.Order{}, err
	}

	return body, nil
}

func (repo *orderRepository) GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error) {
	var body models.CombinedOrderDetails

	query := `
	SELECT orders.id as order_id,orders.final_price,
orders.shipment_status,orders.payment_status,
users.name,users.email,users.phone,
addresses.house_name,addresses.street,
addresses.city,addresses.state,
addresses.pin 
FROM orders INNER JOIN users 
ON orders.user_id = users.id INNER JOIN addresses 
ON orders.address_id = addresses.id WHERE orders.id =?
	`
	if err := repo.DB.Raw(query, orderId).Scan(&body).Error; err != nil {
		err = errors.New("error in getting detailed order through id in repository: " + err.Error())
		return models.CombinedOrderDetails{}, err
	}

	return body, nil
}

func (or *orderRepository) GetFinalPriceOrder(orderID int) (float64, error) {
	var final_price float64
	err := or.DB.Raw("select final_price from orders where id = ?", orderID).Scan(&final_price).Error
	if err != nil {
		return 0.0, errors.New(errmsg.ErrGetDB)
	}
	return final_price, nil
}
func (o *orderRepository) GetItemsByOrderId(orderId int) ([]models.ItemDetails, error) {
	var items []models.ItemDetails

	query := `
	SELECT
    i.product_name,
    oi.quantity,
    i.price,
    oi.total_price
FROM
    orders o
JOIN
    order_items oi ON o.id = oi.order_id
JOIN
    products i ON oi.product_id = i.id
WHERE
    o.id = ?;
	`

	if err := o.DB.Raw(query, orderId).Scan(&items).Error; err != nil {
		return []models.ItemDetails{}, err
	}

	return items, nil
}

func (or *orderRepository) AddTotalToOrder(orderId int, amount float64) error {
	err := or.DB.Raw("update orders set total_amount = ? where id = ? ", amount, orderId).Error
	if err != nil {
		return errors.New(errmsg.ErrWriteDB)
	}
	return nil
}

func (or *orderRepository) PayRazorZero(orderId int) error {
	err := or.DB.Raw("update orders set payment_status = 'PAID' where id = ? ", orderId).Error
	if err != nil {
		return errors.New(errmsg.ErrWriteDB)
	}
	return nil
}
