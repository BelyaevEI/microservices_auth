package consumer

import (
	"context"
)

// Servicer interface for consumer
type Servicer interface {
	RunConsumer(ctx context.Context) error
}
