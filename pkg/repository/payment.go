package repository

import (
	interfaces "WatchHive/pkg/repository/interface"
	"WatchHive/pkg/utils/models"
	"errors"

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

//razor

func (repo *paymentRepository) AddRazorPayDetails(orderId int, razorPayId string) error {
	query := `
	insert into payments (order_id,razer_id) values($1,$2) 
	`
	if err := repo.DB.Exec(query, orderId, razorPayId).Error; err != nil {
		err = errors.New("error in inserting values to razor pay data table" + err.Error())
		return err
	}
	return nil
}

func (pr *paymentRepository) UpdatePaymentDetails(orderId string, paymentId string) error {

	if err := pr.DB.Exec("update payments set payment = $1 where razer_id = $2", paymentId, orderId).Error; err != nil {
		err = errors.New("error in updating the razer pay table " + err.Error())
		return err
	}
	return nil
}

// ------------------------------------------- check payment status ----------------------------------- \\

func (pr *paymentRepository) GetPaymentStatus(orderId string) (bool, error) {
	var paymentStatus string
	err := pr.DB.Raw("select payment_status from orders where id = $1", orderId).Scan(&paymentStatus).Error
	if err != nil {
		return false, err
	}

	// Check if payment status is "PAID"
	isPaid := paymentStatus == "PAID"

	return isPaid, nil
}

func (pr *paymentRepository) UpdatePaymentStatus(status bool, orderId string) error {
	var paymentStatus string
	if status {
		paymentStatus = "PAID"
	} else {
		paymentStatus = "not_paid"
	}

	query := `
		UPDATE orders SET payment_status = $1, shipment_status = 'shipped' WHERE id = $2 
	`
	if err := pr.DB.Exec(query, paymentStatus, orderId).Error; err != nil {
		err = errors.New("error in updating orders payment status: " + err.Error())
		return err
	}
	return nil
}
