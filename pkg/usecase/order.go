package usecase

import (
	"WatchHive/pkg/domain"
	repo_interface "WatchHive/pkg/repository/interface"
	usecase_interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type orderUseCase struct {
	orderRepository   repo_interface.OrderRepository
	cartRepository    repo_interface.CartRepository
	userRepository    repo_interface.UserRepository
	paymentRepository repo_interface.PaymentRepository
}

func NewOrderUseCase(orderRepo repo_interface.OrderRepository, cartRepo repo_interface.CartRepository, userRepo repo_interface.UserRepository, paymentRepo repo_interface.PaymentRepository) usecase_interfaces.OrderUseCase {
	return &orderUseCase{
		orderRepository:   orderRepo,
		cartRepository:    cartRepo,
		userRepository:    userRepo,
		paymentRepository: paymentRepo,
	}

}

func (ou *orderUseCase) Checkout(userID int) (models.CheckoutDetails, error) {
	allUserAddress, err := ou.userRepository.GetAllAddress(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	paymentDetails, err := ou.orderRepository.GetAllPaymentOption()
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	cartItems, err := ou.cartRepository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	grandTotal, err := ou.cartRepository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,
		Total_Price:         grandTotal.FinalPrice,
	}, nil
}

// tyewtuyuiyewiuye

func (ou *orderUseCase) OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (models.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	orderBody.UserID = userID
	cartExist, err := ou.cartRepository.CheckCart(userID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return models.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := ou.userRepository.AddressExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return models.OrderSuccessResponse{}, errors.New("address does not exist")
	}
	PaymentExist, err := ou.paymentRepository.PaymentExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if !PaymentExist {
		return models.OrderSuccessResponse{}, errors.New("paymentmethod does not exist")
	}
	cartItems, err := ou.cartRepository.DisplayCart(orderBody.UserID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	total, err := ou.cartRepository.TotalAmountInCart(orderBody.UserID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	order_id, err := ou.orderRepository.OrderItems(orderBody, total)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if err := ou.orderRepository.AddOrderProducts(order_id, cartItems); err != nil {
		return models.OrderSuccessResponse{}, err
	}

	orderSuccessResponse, err := ou.orderRepository.GetBriefOrderDetails(order_id)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := ou.cartRepository.UpdateCartAfterOrder(userID, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
	}
	return orderSuccessResponse, nil
}

//jhjhsjhsgf

func (ou *orderUseCase) ExecutePurchaseCOD(orderID int) error {
	err := ou.orderRepository.OrderExist(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item  delivered, cannot pay")
	}
	if shipmentStatus == "order placed" {
		return errors.New("item placed, cannot pay")
	}
	if shipmentStatus == "cancelled" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + "so can't paid")
	}
	if shipmentStatus == "processing" {
		return errors.New("the order is already paid")
	}
	err = ou.orderRepository.UpdateOrder(orderID)
	if err != nil {
		return err
	}

	return nil

}

func (or *orderUseCase) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := or.orderRepository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func (or *orderUseCase) CancelOrders(orderID int, userId int) error {
	userTest, err := or.orderRepository.UserOrderRelationship(orderID, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	orderProductDetails, err := or.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := or.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + ", so no point in cancelling")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = or.orderRepository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	err = or.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	return nil

}

func (ou *orderUseCase) GetAllOrdersAdmin(page models.Page) ([]models.CombinedOrderDetails, error) {

	if page.Page == 0 {
		page.Page = 1
	}
	offset := (page.Page - 1) * page.Size

	orderDetail, err := ou.orderRepository.GetAllOrdersAdmin(offset,page.Size)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetail, nil

}
