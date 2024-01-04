package usecase

import (
	"WatchHive/pkg/domain"
	repo_interface "WatchHive/pkg/repository/interface"
	usecase_interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf"
)

type orderUseCase struct {
	orderRepository   repo_interface.OrderRepository
	cartRepository    repo_interface.CartRepository
	userRepository    repo_interface.UserRepository
	paymentRepository repo_interface.PaymentRepository
	walletRepo        repo_interface.WalletRepository
	couponRepo        repo_interface.CouponRepository
}

func NewOrderUseCase(orderRepo repo_interface.OrderRepository, couponRepo repo_interface.CouponRepository, walletRepo repo_interface.WalletRepository, cartRepo repo_interface.CartRepository, userRepo repo_interface.UserRepository, paymentRepo repo_interface.PaymentRepository) usecase_interfaces.OrderUseCase {
	return &orderUseCase{
		orderRepository:   orderRepo,
		cartRepository:    cartRepo,
		userRepository:    userRepo,
		paymentRepository: paymentRepo,
		walletRepo:        walletRepo,
		couponRepo:        couponRepo,
	}

}

func (ou *orderUseCase) Checkout(userID int) (models.CheckoutDetails, error) {
	ok, err := ou.cartRepository.CheckCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	if !ok {
		return models.CheckoutDetails{}, errors.New(errmsg.ErrEmptyCart)
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
		return models.OrderSuccessResponse{}, errors.New(errmsg.ErrEmptyCart)
	}

	addressExist, err := ou.userRepository.AddressExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return models.OrderSuccessResponse{}, errors.New("address " + errmsg.ErrNotExist)
	}
	PaymentExist, err := ou.paymentRepository.PaymentExist(orderBody)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	// var CouponData models.CouponResp
	if orderBody.CouponID > 0 {

		couponExist, err := ou.couponRepo.IsCouponExistByID(orderBody.CouponID)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
		if !couponExist {
			return models.OrderSuccessResponse{}, errors.New(errmsg.ErrCouponExistFalse)
		}

	}

	if !PaymentExist {
		return models.OrderSuccessResponse{}, errors.New("paymentmethod " + errmsg.ErrNotExist)
	}
	cartItems, err := ou.cartRepository.DisplayCart(orderBody.UserID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	total, err := ou.cartRepository.TotalAmountInCart(orderBody.UserID)

	totalOld := total
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	if orderBody.CouponID > 0 {
		couponData, err := ou.couponRepo.GetCouponData(orderBody.CouponID)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
		total -= (total * float64(couponData.OfferPercentage) / 100)
	}

	walletData, err := ou.walletRepo.GetWalletData(orderBody.UserID)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}
	if total < walletData.Amount {

		err := ou.walletRepo.DebitFromWallet(orderBody.UserID, total)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
		total = 0.0
	} else {
		err := ou.walletRepo.DebitFromWallet(orderBody.UserID, walletData.Amount)
		if err != nil {
			return models.OrderSuccessResponse{}, err
		}
		total -= walletData.Amount
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
	var walletDebit models.WalletHistory
	walletDebit.Amount = total
	walletDebit.OrderID = order_id
	walletDebit.Status = "DEBITED"
	walletDebit.WalletID = int(walletData.Id)

	err = ou.walletRepo.AddToWalletHistory(walletDebit)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	orderSuccessResponse, err := ou.orderRepository.GetBriefOrderDetails(order_id)
	if err != nil {
		return models.OrderSuccessResponse{}, err
	}

	orderSuccessResponse.Total = totalOld
	orderSuccessResponse.FinalPrice = total

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
		return errors.New(errmsg.ErrInvalidOId)
	}
	userTest, err := ou.orderRepository.UserOrderRelationship(orderID, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New(errmsg.ErrUserOwnedOrder)
	}
	ok, err := ou.orderRepository.OrderExist(orderID)
	if err != nil {
		return errors.New(errmsg.ErrGetData)
	}
	if !ok {
		return errors.New("order " + errmsg.ErrNotExist)
	}
	orderProductDetails, err := ou.orderRepository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	paymentStatus, err := ou.orderRepository.GetPaymentStatus(orderID)
	if err != nil {
		return err
	}

	if shipmentStatus == "returned" {
		return errors.New(errmsg.ErrReturnedAlready)
	}

	if shipmentStatus == "cancelled" {
		return errors.New(errmsg.ErrCancelAlready)
	}

	if shipmentStatus == "Delivered" {
		return errors.New(errmsg.ErrDeliveredAlready)
	}

	err = ou.orderRepository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	if paymentStatus == "paid" || paymentStatus == "PAID" {
		amount, err := ou.orderRepository.GetFinalPriceOrder(orderID)
		if err != nil {
			return err
		}
		err = ou.walletRepo.AddToWallet(userId, amount)
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
		return errors.New(errmsg.ErrGetData)
	}
	if !ok {
		return errors.New("order" + errmsg.ErrNotExist)
	}
	ShipmentStatus, err := ou.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	paymenType, err := ou.orderRepository.GetPaymentType(orderId)
	if err != nil {
		return err
	}
	paymentStatus, err := ou.orderRepository.GetPaymentStatus(orderId)
	if err != nil {
		return err
	}
	//cod
	if paymenType == 1 {

		if ShipmentStatus == "cancelled" {
			return errors.New(errmsg.ErrCancelAlreadyApprove)
		}
		if ShipmentStatus == "pending" {
			return errors.New(errmsg.ErrPendingApprove)
		}
		if ShipmentStatus == "delivered" {
			return errors.New(errmsg.ErrDeliveredApprove)
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
			return errors.New(errmsg.ErrCancelAlreadyApprove)
		}
		if ShipmentStatus == "pending" {
			return errors.New(errmsg.ErrPendingApprove)
		}
		if ShipmentStatus == "delivered" {
			return errors.New(errmsg.ErrDeliveredApprove)
		}
		if ShipmentStatus == "processing" && paymentStatus == "PAID" {
			err := ou.orderRepository.ApproveOrder(orderId)
			if err != nil {
				return err
			}

			return nil
		}
		if ShipmentStatus == "shipped" {
			// 	err := ou.orderRepository.ApproveRazorPaid(orderId)
			// 	if err != nil {
			// 		return err
			// }
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
		return errors.New(errmsg.ErrInvalidOId)
	}
	ok, err := ou.orderRepository.CheckOrderID(orderId)

	if !ok {
		return errors.New("order " + errmsg.ErrNotExist)
	}
	if err != nil {
		return err
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
		return errors.New(errmsg.ErrCancelAlready)
	}
	if ShipmentStatus == "deliverd" {
		return errors.New(errmsg.ErrDeliveredAlreadyCancel)
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
		return errors.New(errmsg.ErrInvalidData)
	}
	ok, err := ou.orderRepository.CheckOrderID(orderId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("order " + errmsg.ErrNotExist)
	}
	userTest, err := ou.orderRepository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New(errmsg.ErrUserOwnedOrder)
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
		return errors.New(errmsg.ErrCancelAlreadyReturn)
	}
	if shipmentStatus == "pending" {
		return errors.New(errmsg.ErrPendingReturn)
	}
	if shipmentStatus == "processing" {
		return errors.New(errmsg.ErrProcessingReturn)
	}
	if shipmentStatus == "returned" {
		return errors.New(errmsg.ErrReturnedAlready)
	}
	if shipmentStatus == "shipped" {
		return errors.New(errmsg.ErrShippedReturn)
	}
	amount, err := ou.orderRepository.GetFinalPriceOrder(orderId)
	if err != nil {
		return err
	}
	if paymenType == 1 {

		if shipmentStatus == "delivered" {

			err = ou.orderRepository.ReturnOrderCod(orderId)
			if err != nil {
				return err
			}
			err = ou.walletRepo.AddToWallet(userId, amount)
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
			err = ou.walletRepo.AddToWallet(userId, amount)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (or *orderUseCase) PrintInvoice(orderId, userId int) (*gofpdf.Fpdf, error) {
	if orderId < 0 {
		return nil, errors.New(errmsg.ErrInvalidData)
	}
	ok, err := or.orderRepository.CheckOrderID(orderId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("order " + errmsg.ErrNotExist)
	}
	userTest, err := or.orderRepository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return nil, err
	}
	if userTest != userId {
		return nil, errors.New(errmsg.ErrUserOwnedOrder)
	}
	if orderId < 1 {
		return nil, errors.New("enter a valid order id")
	}

	order, err := or.orderRepository.GetDetailedOrderThroughId(orderId)
	if err != nil {
		return nil, err
	}

	items, err := or.orderRepository.GetItemsByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	if order.ShipmentStatus != "delivered" {
		return nil, errors.New(errmsg.ErrDeliverInvoice)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 30)
	pdf.SetTextColor(31, 73, 125)
	pdf.Cell(0, 20, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "I", 14)
	pdf.SetTextColor(51, 51, 51)
	pdf.Cell(0, 10, "Customer Details")
	pdf.Ln(10)
	customerDetails := []string{
		"Name: " + order.Name,
		"House Name: " + order.HouseName,
		"Street: " + order.Street,
		"State: " + order.State,
		"City: " + order.City,
	}
	for _, detail := range customerDetails {
		pdf.Cell(0, 10, detail)
		pdf.Ln(10)
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, "Item", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Final Price", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(255, 255, 255)
	for _, item := range items {
		pdf.CellFormat(40, 10, item.ProductName, "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.Price, 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, strconv.Itoa(item.Quantity), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.Price*float64(item.Quantity), 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.Ln(10)
	}
	pdf.Ln(10)

	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Total Price:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(totalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	offerApplied := totalPrice - order.FinalPrice

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Offer Applied:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(offerApplied, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Final Amount:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(order.FinalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 12)
	pdf.Cell(0, 10, "Generated by Watch Hive India Pvt Ltd. - "+time.Now().Format("2006-01-02 15:04:05"))
	pdf.Ln(10)

	return pdf, nil
}
