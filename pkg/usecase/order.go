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
	walletRepo repo_interface.WalletRepository
}

func NewOrderUseCase(orderRepo repo_interface.OrderRepository,walletRepo repo_interface.WalletRepository, cartRepo repo_interface.CartRepository, userRepo repo_interface.UserRepository, paymentRepo repo_interface.PaymentRepository) usecase_interfaces.OrderUseCase {
	return &orderUseCase{
		orderRepository:   orderRepo,
		cartRepository:    cartRepo,
		userRepository:    userRepo,
		paymentRepository: paymentRepo,
		walletRepo: walletRepo,
	}

}

func (ou *orderUseCase) Checkout(userID int) (models.CheckoutDetails, error) {
	ok, err := ou.cartRepository.CheckCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, errors.New("error in getting details")
	}
	if !ok {
		return models.CheckoutDetails{}, errors.New("no items in cart")
	}
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

	// here placing order

	err = ou.orderRepository.UpdateOrder(order_id)
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

	orderSuccessResponse, err := ou.orderRepository.GetBriefOrderDetails(order_id)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	return orderSuccessResponse, nil
}

func (ou *orderUseCase) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := ou.orderRepository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func (ou *orderUseCase) CancelOrders(orderID int, userId int) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}
	userTest, err := ou.orderRepository.UserOrderRelationship(orderID, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	ok, err := ou.orderRepository.OrderExist(orderID)
	if err != nil {
		return errors.New("error in getting data")
	}
	if !ok {
		return errors.New("order is not exist")
	}
	orderProductDetails, err := ou.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
paymentStatus,err := ou.orderRepository.GetPaymentStatus(orderID)
if err != nil {
	return err
}

	if shipmentStatus == "pending" || shipmentStatus == "returned" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + ", so no point in cancelling")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}

	if shipmentStatus == "Delivered" {
		return errors.New("the order is delivered, you can return it")
	}

	err = ou.orderRepository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	if paymentStatus == "paid" || paymentStatus == "PAID" {
		amount,err := ou.orderRepository.GetFinalPriceOrder(orderID)
		if err!= nil{
			return err
		}
		err = ou.walletRepo.AddToWallet(userId,amount)
		if err != nil {
			return err
		}
	}
	err = ou.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
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

	orderDetail, err := ou.orderRepository.GetAllOrdersAdmin(offset, page.Size)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetail, nil

}

func (ou *orderUseCase) ApproveOrder(orderId int) error {
	if orderId <= 0 {
		return errors.New("invalid id")
	}
	ok, err := ou.orderRepository.OrderExist(orderId)
	if err != nil {
		return errors.New("error in getting data")
	}
	if !ok {
		return errors.New("order is not exist")
	}
	ShipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	paymenType, err := ou.orderRepository.GetPaymentType(orderId)
	if err != nil {
		return err
	}
//cod
	if paymenType == 1 {

		if ShipmentStatus == "cancelled" {
			return errors.New("the order is cancelled,cannot approve it")
		}
		if ShipmentStatus == "pending" {
			return errors.New("the order is pending,cannot approve it")
		}
		if ShipmentStatus == "delivered" {
			return errors.New("the order is already deliverd")
		}
		if ShipmentStatus == "processing" {
			err := ou.orderRepository.ApproveOrder(orderId)
			if err != nil {
				return err
			}

			return nil
		}

		if ShipmentStatus == "shipped" {
			err := ou.orderRepository.ApproveCodPaid(orderId)
			if err != nil {
				return err
			}

			return nil
		}

		if ShipmentStatus == "returned" {
			err := ou.orderRepository.ApproveCodReturn(orderId)
			if err != nil {
				return err
			}
		}
	}
	// razorpay
	if paymenType == 2 {
		if ShipmentStatus == "cancelled" {
			return errors.New("the order is cancelled,cannot approve it")
		}
		if ShipmentStatus == "pending" {
			return errors.New("the order is pending,cannot approve it")
		}
		if ShipmentStatus == "delivered" {
			return errors.New("the order is already deliverd")
		}
		if ShipmentStatus == "processing" {
			err := ou.orderRepository.ApproveOrder(orderId)
			if err != nil {
				return err
			}

			return nil
		}
		if ShipmentStatus == "shipped" {
			err := ou.orderRepository.ApproveRazorPaid(orderId)
			if err != nil {
				return err
			}
			err = ou.orderRepository.ApproveRazorDelivered(orderId)
			if err != nil {
				return err
			}

			return nil

		}
	}

	// if the shipment status is not processing or cancelled. Then it is defenetely cancelled
	return nil
}

func (ou *orderUseCase) CancelOrderFromAdmin(orderId int) error {
	if orderId <= 0 {
		return errors.New("invalid order id")
	}
	ok, err := ou.orderRepository.CheckOrderID(orderId)
	fmt.Println(err)
	if !ok {
		return errors.New("order does not exist")
	}
	orderProduct, err := ou.orderRepository.GetProductDetailsFromOrders(orderId)
	if err != nil {
		return err
	}

	ShipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	if ShipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled")
	}
	if ShipmentStatus == "deliverd" {
		return errors.New("the order is delivered cannot be cancelled")
	}
	err = ou.orderRepository.CancelOrders(orderId)
	if err != nil {
		return err
	}
	err = ou.orderRepository.UpdateStockOfProduct(orderProduct)
	if err != nil {
		return err
	}
	return nil
}

func (ou *orderUseCase) ReturnOrder(orderId, userId int) error {

	if orderId < 0 {
		return errors.New("invalid order id")
	}
	ok, err := ou.orderRepository.CheckOrderID(orderId)
	fmt.Println(err)
	if !ok {
		return errors.New("order does not exist")
	}
	userTest, err := ou.orderRepository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}

	shipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	paymenType, err := ou.orderRepository.GetPaymentType(orderId)
	if err != nil {
		return err
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is cancelled,cannot return it")
	}
	if shipmentStatus == "pending" {
		return errors.New("the order is pending,cannot return it")
	}
	if shipmentStatus == "processing" {
		return errors.New("the order is processing cannot return it")
	}
	if shipmentStatus == "returned" {
		return errors.New("the order is already returned")
	}
	if shipmentStatus == "shipped" {
		return errors.New("the order is shipped cannot return it")
	}

	if paymenType == 1 {

		if shipmentStatus == "delivered" {

			err = ou.orderRepository.ReturnOrderCod(orderId)
			if err != nil {
				return err
			}

		}
	}
	if paymenType == 2 {
		if shipmentStatus == "delivered" {
		err = ou.orderRepository.ReturnOrderRazorPay(orderId)
		if err != nil {
			return err
		}
	}
	}
	return nil
}
