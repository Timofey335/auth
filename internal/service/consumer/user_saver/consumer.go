package user_saver

import (
	"context"

	"github.com/Timofey335/platform_common/pkg/kafka"

	"github.com/Timofey335/auth/internal/repository"
	def "github.com/Timofey335/auth/internal/service"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
	consumer       kafka.Consumer
}

// NewService
func NewService(
	userRepository repository.UserRepository,
	consumer kafka.Consumer,
) *service {
	return &service{
		userRepository: userRepository,
		consumer:       consumer,
	}
}

// RunConsumer - запускает consumer
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

		errChan <- s.consumer.Consume(ctx, "create_user", s.UserSaveHandler)
	}()

	return errChan
}
