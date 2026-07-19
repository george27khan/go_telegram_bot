package state

import "go_telegram_bot/internal/domain/entity"

type StateUseCase struct {
	Current entity.UserState
	In      entity.UserState
}

func NewStateUseCase() *StateUseCase {
	return &StateUseCase{}
}

func (s *StateUseCase) Next(current entity.UserState, in entity.UserState) entity.UserState {
	return ""
}
