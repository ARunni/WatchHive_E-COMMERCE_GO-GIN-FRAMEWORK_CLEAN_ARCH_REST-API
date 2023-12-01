package interfaces

import "WatchHive/pkg/utils/models"

type PaymentRepository interface {
	PaymentExist(orderBody models.OrderIncoming) (bool, error)
	PaymentMethodID(orderID int) (int, error) 
	CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	AddPaymentMethod(pay models.NewPaymentMethod) (models.PaymentDetails, error)
}

