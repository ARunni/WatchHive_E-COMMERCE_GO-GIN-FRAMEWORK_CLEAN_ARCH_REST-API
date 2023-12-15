package interfaces

import (
	"WatchHive/pkg/utils/models"
)

type PaymentUseCase interface {
	PaymentMethodID(order_id int) (int, error)
	AddPaymentMethod(payment models.NewPaymentMethod) (models.PaymentDetails, error)
	MakePaymentRazorpay(orderId, userId int) (models.CombinedOrderDetails, string, error) 
	SavePaymentDetails(paymentId, razorId, orderId string) error 
	
}
