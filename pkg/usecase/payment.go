package usecase

import (
	"WatchHive/pkg/config"
	interfaces_repo "WatchHive/pkg/repository/interface"
	interfaces_usecase "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"errors"

	"github.com/razorpay/razorpay-go"
)

type paymentUseCase struct {
	paymentRepository interfaces_repo.PaymentRepository
	orderRepo         interfaces_repo.OrderRepository
	cfg               config.Config
}

func NewPaymentUseCase(repo interfaces_repo.PaymentRepository, orderRepo interfaces_repo.OrderRepository, cfg config.Config) interfaces_usecase.PaymentUseCase {
	return &paymentUseCase{
		paymentRepository: repo,
		orderRepo:         orderRepo,
		cfg:               cfg,
	}

}

func (pu *paymentUseCase) PaymentMethodID(order_id int) (int, error) {
	if order_id <= 0 {
		return 0, errors.New(errmsg.ErrInvalidOId)
	}

	id, err := pu.paymentRepository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (pu *paymentUseCase) AddPaymentMethod(payment models.NewPaymentMethod) (models.PaymentDetails, error) {
	if payment.PaymentName == "" {
		return models.PaymentDetails{}, errors.New("payment method" + errmsg.ErrFieldEmpty)
	}
	exists, err := pu.paymentRepository.CheckIfPaymentMethodAlreadyExists(payment.PaymentName)
	if err != nil {
		return models.PaymentDetails{}, err
	}
	if exists {
		return models.PaymentDetails{}, errors.New("payment method " + errmsg.ErrExistTrue)
	}
	paymentadd, err := pu.paymentRepository.AddPaymentMethod(payment)
	if err != nil {
		return models.PaymentDetails{}, err
	}
	return paymentadd, nil
}

//razor

func (pu *paymentUseCase) MakePaymentRazorpay(orderId, userId int) (models.CombinedOrderDetails, string, error) {

	if orderId <= 0 || userId <= 0 {
		return models.CombinedOrderDetails{}, "", errors.New(errmsg.ErrInvalidData)
	}

	order, err := pu.orderRepo.GetOrder(orderId)
	if err != nil {
		err = errors.New(errmsg.ErrGetData + err.Error())
		return models.CombinedOrderDetails{}, "", err
	}

	client := razorpay.NewClient(pu.cfg.RazorPay_key_id, pu.cfg.RazorPay_key_secret)

	data := map[string]interface{}{
		"amount":   int(order.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return models.CombinedOrderDetails{}, "", nil
	}

	razorPayOrderId := body["id"].(string)

	err = pu.paymentRepository.AddRazorPayDetails(orderId, razorPayOrderId)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}
	body2, err := pu.orderRepo.GetDetailedOrderThroughId(int(order.ID))
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}

	return body2, razorPayOrderId, nil
}

func (pu *paymentUseCase) SavePaymentDetails(paymentId, razorId, orderId string) error {

	status, err := pu.paymentRepository.GetPaymentStatus(orderId)
	if err != nil {
		return err
	}

	if !status {
		err = pu.paymentRepository.UpdatePaymentDetails(razorId, paymentId)
		if err != nil {
			return err
		}

		err = pu.paymentRepository.UpdatePaymentStatus(true, orderId)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(errmsg.ErrAlreadyPaid)

}
