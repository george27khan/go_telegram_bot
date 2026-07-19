package setting

import (
	"go_telegram_bot/internal/domain/entity"
)

type SettingRepository interface {
}

type SettingUseCase struct {
	SettingRepository SettingRepository
}

func NewSettingUseCase(settingRepository SettingRepository) *SettingUseCase {
	return &SettingUseCase{
		SettingRepository: settingRepository,
	}
}

func (s *SettingUseCase) Next(current entity.UserState, in entity.UserState) entity.UserState {
	return ""
}
