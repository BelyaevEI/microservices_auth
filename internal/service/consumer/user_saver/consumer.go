package usersaver

import (
	"context"

	"github.com/BelyaevEI/microservices_auth/internal/client/kafka"
	"github.com/BelyaevEI/microservices_auth/internal/repository"
	def "github.com/BelyaevEI/microservices_auth/internal/service/consumer"
)

var _ def.Servicer = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	consumer       kafka.Consumer
}

// NewService create instance for consumer
func NewService(
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		userRepository: userRepository,
		consumer:       consumer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *service) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, "test-topic", s.UserSaveHandler)
	}()

	return errChan
}
