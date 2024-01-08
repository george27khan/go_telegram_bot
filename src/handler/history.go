package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	emp "go_telegram_bot/database/employee"
	sched "go_telegram_bot/database/schedule"
)

// entryHandler функция вывода меню просмотра записей
func entryHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	kb := inline.New(b).
		Row().
		Button("Текущая запись", []byte(""), actualEntryHandler).
		Row().
		Button("История записей", []byte(""), histEntryHandler).
		Row().
		Button("⬅️ Назад", []byte(""), BackStartHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// actualEntryHandler функция вывода актуальной записи
func actualEntryHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	schedule, err := sched.GetByUser(ctx, mes.Chat.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе поиска записи " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	FIO, err := emp.GetFIO(ctx, schedule.IdEmployee)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе поиска сотрудника " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Вы записаны к " + FIO + " на " + schedule.VisitDt.Format(datetimeFormat)),
		ReplyMarkup: inline.New(b).Button("⬅️ Назад", []byte(""), entryHandler),
		ParseMode:   models.ParseModeHTML,
	})
}

// histEntryHandler функция вывод истории записей
func histEntryHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	var (
		res string
	)
	schedSlice, err := sched.GetAllByUser(ctx, mes.Chat.ID)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе поиска записей " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	for _, schedule := range schedSlice {
		FIO, err := emp.GetFIO(ctx, schedule.IdEmployee)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    mes.Chat.ID,
				Text:      highlightTxt("Ошибка в процессе поиска сотрудника " + err.Error()),
				ParseMode: models.ParseModeHTML,
			})
			return
		} else {
			res = res + schedule.VisitDt.Format(datetimeFormat) + " " + FIO + "\n"
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt(res),
		ReplyMarkup: inline.New(b).Button("⬅️ Назад", []byte(""), entryHandler),
		ParseMode:   models.ParseModeHTML,
	})
}
