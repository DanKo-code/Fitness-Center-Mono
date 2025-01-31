package usecase

type StripeUseCase interface {
	ArchiveStripeProduct(stripePriceId string) error
	CreateStripeProductAndPrice(name string, amount int64, currency string) (string, error)
	CreateStripePriceAndAssignToProductDeactivateOldPrices(stripePriceId string, amount int64, currency string) (string, error)
	UpdateStripeProductName(stripePriceId string, newName string) error
}
