package consumer

import (
	"context"
)

// ConsumerService interface for consumer
type Servicer interface {
	RunConsumer(ctx context.Context) error
}
