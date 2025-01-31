package background

import (
	"context"
	"time"
)

type OrderExpiredChecker interface {
	StartOrderExpiredChecker(ctx context.Context, interval time.Duration, stopChan <-chan struct{})
}
