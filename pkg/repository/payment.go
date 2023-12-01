package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"

	"gorm.io/gorm"
)

type paymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(DB *gorm.DB) interfaces.PaymentRepository {
	return &paymentRepository{
		DB: DB,
	}

}

func (pr *paymentRepository) PaymentExist(orderBody models.OrderIncoming) (bool, error) {
	var count int
	if err := pr.DB.Raw("SELECT count(*) FROM payment_methods WHERE id = ?", orderBody.PaymentID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}

func (pr *paymentRepository) PaymentMethodID(orderID int) (int, error) {
	var a int
	err := pr.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil
}

func (pr *paymentRepository) AddPaymentMethod(pay models.NewPaymentMethod) (models.PaymentDetails, error) {
	var payment string
	if err := pr.DB.Raw("INSERT INTO payment_methods (payment_name) VALUES (?) RETURNING payment_name", pay.PaymentName).Scan(&payment).Error; err != nil {
		return models.PaymentDetails{}, err
	}
	var paymentResponse models.PaymentDetails
	err := pr.DB.Raw("SELECT id, payment_name FROM payment_methods WHERE payment_name = ?", payment).Scan(&paymentResponse).Error
	if err != nil {
		return models.PaymentDetails{}, err
	}
	return paymentResponse, nil

}

func (pr *paymentRepository) CheckIfPaymentMethodAlreadyExists(payment string) (bool, error) {
	var count int64
	err := pr.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE payment_name = $1", payment).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
