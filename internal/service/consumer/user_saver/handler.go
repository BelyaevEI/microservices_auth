package usersaver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/IBM/sarama"
)

func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.UserCreate{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	log.Printf("Note with id %d created\n", id)

	return nil
}
