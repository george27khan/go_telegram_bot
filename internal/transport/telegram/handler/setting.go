package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"io"
	"log/slog"
	"net/http"
)

type SettingUseCase interface {
	//GetCalendarDays(ctx context.Context) (dateFrom time.Time, dateTo time.Time, err error)
	//GetDayTimes(ctx context.Context, date time.Time) (timeFrom time.Time, timeTo time.Time, err error)
	//GetSessionTimeMinute(ctx context.Context) (time.Duration, error)
	//IsFreeTime(ctx context.Context, visitTime time.Time) (ok bool, err error)
	//GetFreeEmployee(ctx context.Context, visitTime time.Time) ([]entity.Employee, error)
}

type SettingHandler struct {
	SettingUseCase SettingUseCase
	logger         *slog.Logger
}

func NewSettingHandler(settingUseCase SettingUseCase, logger *slog.Logger) SettingHandler {
	return SettingHandler{
		SettingUseCase: settingUseCase,
		logger:         logger,
	}
}

// GetSetting функция вывода меню настроек
func (s *SettingHandler) GetSetting(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	kb := inline.New(b).
		Row().
		Button("Должность", []byte(""), nil). //positionSettingHandler
		Row().
		Button("Сотрудник", []byte(""), nil). //empSettingHandler
		Row().
		Button("Заполнить данные", []byte(""), nil). //empInitHandler
		Row().
		Button("⬅️ Назад", []byte(""), BackStartHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// empInitHandler функция для заполнения тестовыми данными
//func empInitHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
//	if err := fill.AddPositions(ctx); err != nil {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID:    mes.Chat.ID,
//			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
//			ParseMode: models.ParseModeHTML,
//		})
//		return
//	}
//
//	if err := fill.AddEmployees(ctx); err != nil {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID:    mes.Chat.ID,
//			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
//			ParseMode: models.ParseModeHTML,
//		})
//		return
//	}
//	b.SendMessage(ctx, &bot.SendMessageParams{
//		ChatID:    mes.Chat.ID,
//		Text:      highlightTxt("Таблицы успешно заполнены"),
//		ParseMode: models.ParseModeHTML,
//	})
//}

func downloadFile(URL string) ([]byte, error) {
	//Get the response bytes from the url
	var bytes []byte
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		bytes, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

// BackSettingHandler функция возврата в меню настроек
func (s *SettingHandler) BackSettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	s.GetSetting(ctx, b, mes, data)
}
