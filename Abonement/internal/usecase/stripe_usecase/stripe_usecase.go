package stripe_usecase

import (
	"fmt"
	"github.com/DanKo-code/Fitness-Center-Abonement/pkg/logger"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"log"
)

type StripeUseCase struct {
	stripeKey string
}

func NewStripeUseCase(stripeKey string) *StripeUseCase {
	return &StripeUseCase{
		stripeKey: stripeKey,
	}
}

func (suc *StripeUseCase) ArchiveStripeProduct(stripePriceId string) error {
	stripe.Key = suc.stripeKey

	priceObject, err := price.Get(stripePriceId, nil)
	if err != nil {
		logger.ErrorLogger.Printf("Error getting stripe id: %v\n", err)
		return err
	}

	productID := priceObject.Product.ID

	params := &stripe.ProductParams{
		Active: stripe.Bool(false),
	}

	updatedProduct, err := product.Update(productID, params)
	if err != nil {
		logger.ErrorLogger.Printf("Ошибка при архивировании продукта: %v\n", err)
		return err
	}

	logger.InfoLogger.Printf("Product %s successfully has been archived\n", updatedProduct.ID)
	return nil
}

func (suc *StripeUseCase) CreateStripeProductAndPrice(name string, amount int64, currency string) (string, error) {
	stripe.Key = suc.stripeKey

	productParams := &stripe.ProductParams{
		Name: stripe.String(name),
	}
	createdProduct, err := product.New(productParams)
	if err != nil {
		return "", fmt.Errorf("failed to create stripe product: %v", err)
	}
	logger.InfoLogger.Printf("stripe product has been created: %s\n", createdProduct.ID)

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(createdProduct.ID),
		UnitAmount: stripe.Int64(amount),
		Currency:   stripe.String(currency),
	}

	createdPrice, err := price.New(priceParams)
	if err != nil {
		return "", fmt.Errorf("failed to create stripe price: %v", err)
	}
	fmt.Printf("stripe price has been created: %s\n", createdPrice.ID)

	return createdPrice.ID, nil
}

func (suc *StripeUseCase) CreateStripePriceAndAssignToProductDeactivateOldPrices(stripePriceId string, amount int64, currency string) (string, error) {
	stripe.Key = suc.stripeKey

	priceObject, err := price.Get(stripePriceId, nil)
	if err != nil {
		logger.ErrorLogger.Printf("Error getting stripe id: %v\n", err)
		return "", err
	}

	productID := priceObject.Product.ID

	newPriceParams := &stripe.PriceParams{
		UnitAmount: stripe.Int64(amount),
		Currency:   stripe.String(currency),
		Product:    stripe.String(productID),
	}

	newPrice, err := price.New(newPriceParams)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to create new price: %v", err)
		return "", fmt.Errorf("failed to create stripe price: %v", err)
	}

	err = deactivateOldPrices(productID, newPrice.ID)
	if err != nil {
		return "", err
	}

	return newPrice.ID, nil
}

func (suc *StripeUseCase) UpdateStripeProductName(stripePriceId string, newName string) error {
	stripe.Key = suc.stripeKey

	priceObject, err := price.Get(stripePriceId, nil)
	if err != nil {
		logger.ErrorLogger.Printf("Error getting stripe id: %v\n", err)
		return err
	}

	productID := priceObject.Product.ID

	productParams := &stripe.ProductParams{
		Name: stripe.String(newName),
	}

	updatedProduct, err := product.Update(productID, productParams)
	if err != nil {
		return fmt.Errorf("failed to update stripe product name: %v", err)
	}

	logger.InfoLogger.Printf("stripe product name has been updated: %s -> %s\n", updatedProduct.ID, newName)
	return nil
}

func deactivateOldPrices(productID string, newPrice string) error {
	params := &stripe.PriceListParams{
		Product: stripe.String(productID),
		Active:  stripe.Bool(true),
	}

	iter := price.List(params)
	for iter.Next() {
		p := iter.Price()
		if p.ID != newPrice {
			_, err := price.Update(p.ID, &stripe.PriceParams{
				Active: stripe.Bool(false),
			})
			if err != nil {
				log.Printf("Failed to deactivate price %s: %v", p.ID, err)
				continue
			}
			fmt.Printf("Deactivated price: %s\n", p.ID)
		}
	}
	if err := iter.Err(); err != nil {
		logger.ErrorLogger.Printf("Error listing prices: %v", err)
		return err
	}

	return nil
}
