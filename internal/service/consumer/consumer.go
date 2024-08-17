package consumer

import (
	"context"
)

// ConsumerService interface for consumer
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
