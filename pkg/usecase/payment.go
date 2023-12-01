package usecase

import (
	interfaces_repo "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
)

type paymentUseCase struct {
	paymentRepository interfaces_repo.PaymentRepository
}

func NewPaymentUseCase(repo interfaces_repo.PaymentRepository) interfaces.PaymentUseCase {
	return &paymentUseCase{
		paymentRepository: repo,
	}

}

func (pu *paymentUseCase) PaymentMethodID(order_id int) (int, error) {
	id, err := pu.paymentRepository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (pu *paymentUseCase) AddPaymentMethod(payment models.NewPaymentMethod) (models.PaymentDetails, error) {
	exists, err := pu.paymentRepository.CheckIfPaymentMethodAlreadyExists(payment.PaymentName)
	if err != nil {
		return models.PaymentDetails{}, err
	}
	if exists {
		return models.PaymentDetails{}, errors.New("payment method already exists")
	}
	paymentadd, err := pu.paymentRepository.AddPaymentMethod(payment)
	if err != nil {
		return models.PaymentDetails{}, err
	}
	return paymentadd, nil
}
