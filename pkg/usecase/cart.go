package usecase

import (
	interfaces_repo "WatchHive/pkg/repository/interface"
	interfaces "WatchHive/pkg/usecase/interface"
	"WatchHive/pkg/utils/models"
	"errors"
)

type cartUseCase struct {
	cartRepository    interfaces_repo.CartRepository
	productRepository interfaces_repo.ProductRepository
}

func NewCartUseCase(repoc interfaces_repo.CartRepository, repop interfaces_repo.ProductRepository) interfaces.CartUseCase {
	return &cartUseCase{
		cartRepository:    repoc,
		productRepository: repop,
	}
}

func (cu *cartUseCase) AddToCart(cart models.AddCart) (models.CartResponse, error) {

	if cart.ProductID < 1 || cart.UserID < 1 {
		return models.CartResponse{}, errors.New("invalid product id or user id ")
	}
	if cart.Quantity < 1 {
		return models.CartResponse{}, errors.New("quantity must be 1 or greater")
	}
	is_available, err := cu.productRepository.CheckProductAvailable(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if !is_available {
		return models.CartResponse{}, errors.New("product is not available")
	}
	stock, err := cu.cartRepository.CheckStock(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if stock < int(cart.Quantity) {
		return models.CartResponse{}, errors.New("out of stock")
	}

	price, err := cu.productRepository.GetPriceOfProduct(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	QuantityOfProductInCart, err := cu.cartRepository.QuantityOfProductInCart(cart.UserID, int(cart.ProductID))
	if err != nil {

		return models.CartResponse{}, err
	}
	if (QuantityOfProductInCart + int(cart.Quantity)) > 20 {
		return models.CartResponse{}, errors.New("limit exceeds")
	}

	finalPrice := (price * float64(cart.Quantity))

	if QuantityOfProductInCart == 0 {
		err := cu.cartRepository.AddToCart(cart.UserID, int(cart.ProductID), int(cart.Quantity), finalPrice)
		if err != nil {

			return models.CartResponse{}, err
		}

	} else {
		currentTotal, err := cu.cartRepository.TotalPriceForProductInCart(cart.UserID, int(cart.ProductID))
		if err != nil {
			return models.CartResponse{}, err
		}
		err = cu.cartRepository.UpdateCart(QuantityOfProductInCart+int(cart.Quantity), currentTotal+finalPrice, cart.UserID, int(cart.ProductID))
		if err != nil {
			return models.CartResponse{}, err
		}
	}

	cartDetails, err := cu.cartRepository.DisplayCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.cartRepository.GetTotalPrice(cart.UserID)
	if err != nil {

		return models.CartResponse{}, err
	}

	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil
}

func (cu *cartUseCase) ListCartItems(userID int) (models.CartResponse, error) {
	cartDetails, err := cu.cartRepository.DisplayCart(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.cartRepository.GetTotalPrice(userID)
	if err != nil {

		return models.CartResponse{}, err
	}

	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil
}

func (cu *cartUseCase) UpdateProductQuantityCart(cart models.AddCart) (models.CartResponse, error) {

	ok, err := cu.cartRepository.CheckCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, errors.New("error in checking cart")
	}
	if !ok {
		return models.CartResponse{}, errors.New("cart is empty")
	}
	if cart.Quantity < 1 || cart.ProductID < 1 {
		return models.CartResponse{}, errors.New("invalid product id or quantity")
	}
	is_available, err := cu.productRepository.CheckProductAvailable(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}
	if !is_available {
		return models.CartResponse{}, errors.New("product is not available")
	}
	stock, err := cu.cartRepository.CheckStock(int(cart.ProductID))
	if err != nil {
		return models.CartResponse{}, err
	}

	if stock < int(cart.Quantity) {
		return models.CartResponse{}, errors.New("out of stock")
	}

	if int(cart.Quantity) > 20 {
		return models.CartResponse{}, errors.New("limit exceeds")
	}

	err = cu.cartRepository.UpdateProductQuantityCart(cart)
	if err != nil {
		return models.CartResponse{}, err
	}

	cartDetails, err := cu.cartRepository.DisplayCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.cartRepository.GetTotalPrice(cart.UserID)
	if err != nil {

		return models.CartResponse{}, err
	}

	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}

func (cu *cartUseCase) RemoveFromCart(cart models.RemoveFromCart) (models.CartResponse, error) {

	if cart.ProductID < 1 {
		return models.CartResponse{}, errors.New("product id cannot be empty")
	}

	is_available, err := cu.cartRepository.CheckCart(cart.UserID)
	if !is_available {
		return models.CartResponse{}, err
	}

	err = cu.cartRepository.RemoveFromCart(cart)
	if err != nil {
		return models.CartResponse{}, err
	}

	cartDetails, err := cu.cartRepository.DisplayCart(cart.UserID)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cu.cartRepository.GetTotalPrice(cart.UserID)
	if err != nil {

		return models.CartResponse{}, err
	}

	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}
