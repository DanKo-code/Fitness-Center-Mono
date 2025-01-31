package order_background

import (
	"context"
	"github.com/DanKo-code/FitnessCenter-Order/internal/usecase"
	"github.com/DanKo-code/FitnessCenter-Order/pkg/logger"
	"time"
)

type OrderExpiredChecker struct {
	usecase usecase.OrderUseCase
}

func NewOrderExpiredChecker(useCase usecase.OrderUseCase) *OrderExpiredChecker {
	return &OrderExpiredChecker{
		usecase: useCase,
	}
}

func (aec *OrderExpiredChecker) StartOrderExpiredChecker(ctx context.Context, interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := aec.usecase.SetExpiredOrdersTasks(ctx); err != nil {
				logger.ErrorLogger.Printf("Error updating expired orders: %v", err)
				return
			}
			logger.InfoLogger.Println("Updated expired orders")
		case <-stopChan:
			logger.InfoLogger.Println("Stopping expired orders checker")
			return
		}
	}
}
