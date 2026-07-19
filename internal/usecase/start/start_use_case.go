package start

import (
	"context"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/transport/telegram/handler"
)

var _ handler.StartUseCase = (*StartUseCase)(nil)

type ClientRepository interface {
	IsExists(ctx context.Context, clientID entity.ClientID) (bool, error)
	Create(ctx context.Context, client entity.Client) error
}

type StartUseCase struct {
	ClientRepository ClientRepository
}

func NewStartUseCase(clientRepository ClientRepository) *StartUseCase {
	return &StartUseCase{ClientRepository: clientRepository}
}

func (su *StartUseCase) Start(ctx context.Context, client entity.Client) error {
	exists, err := su.ClientRepository.IsExists(ctx, client.ID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if err = su.ClientRepository.Create(ctx, client); err != nil {
		return err
	}
	return nil
}
