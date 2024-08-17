package kafka

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/client/kafka/consumer"
)

// Consumer simple interface
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
