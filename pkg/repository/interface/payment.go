package interfaces

import "WatchHive/pkg/utils/models"

type PaymentRepository interface {
	PaymentExist(orderBody models.OrderIncoming) (bool, error)
	PaymentMethodID(orderID int) (int, error) 
	CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	AddPaymentMethod(pay models.NewPaymentMethod) (models.PaymentDetails, error)
	AddRazorPayDetails(orderId int, razorPayId string) error 
	UpdatePaymentDetails(orderId string, paymentId string) error 
	GetPaymentStatus(orderId string) (bool, error)
	UpdatePaymentStatus(status bool, orderId string) error 
}


