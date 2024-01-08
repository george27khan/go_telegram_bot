package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	pstn "go_telegram_bot/database/position"
	"go_telegram_bot/src/state"
	"strconv"
	"strings"
)

// positionSettingHandler функция вывода настроек должностей
func positionSettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b).
		Row().
		Button("Добавить должность", []byte(""), positionAddHandler).
		Row().
		Button("Удалить должность", []byte(""), positionDelHandler).
		Row().
		Button("Вывести должности", []byte(""), positionShowHandler).
		Row().
		Button("⬅️Назад", []byte(""), BackSettingHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// positionNameKB функция формирования клавиатуры из должностей
func positionNameKB(ctx context.Context, b *bot.Bot, action func(context.Context, *bot.Bot, *models.Message, []byte)) *inline.Keyboard {
	kb := inline.New(b)
	positions, err := pstn.SelectAll(ctx)
	if err != nil {
		fmt.Println("PositionNameKB error:", err)
	}
	for _, position := range positions {
		kb = kb.Row().Button(position.PositionName, []byte(strconv.Itoa(position.Id)), action)
	}
	return kb
}

// positionDelHandler функция вывода позиции для удаления
func positionDelHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите должность на удаления"),
		ReplyMarkup: positionNameKB(ctx, b, delPosition).Row().Button("Назад", []byte(""), positionSettingHandler),
		ParseMode:   models.ParseModeHTML,
	})
}

// delPosition функция для удаления выбранной должности
func delPosition(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var (
		res string
	)
	id, err := strconv.Atoi(string(data))
	if err != nil {
		res = "❌ Ошибка в процессе удаления позиции: " + err.Error()
	}
	if err := pstn.DeleteById(ctx, id); err != nil {
		res = "❌ Ошибка в процессе удаления позиции: " + err.Error()
	} else {
		res = "✅ Должность удалена"
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt(res),
		ParseMode: models.ParseModeHTML,
	})

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: menuKeyboard(b),
		ParseMode:   models.ParseModeHTML,
	})
}

// positionAddHandler функция запроса создания новой должности
func positionAddHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	state.Set(ctx, "user_state", strconv.FormatInt(mes.Chat.ID, 10), "PositionAddHandler")
	positions, _ := pstn.SelectAllStr(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Введите новое название должности НЕ из списка:" + "\n" + strings.Join(positions, "\n")),
		ReplyMarkup: inline.New(b).Button("Назад", []byte(""), positionSettingHandler),
		ParseMode:   models.ParseModeHTML,
	})
}

// addPosition функция для сохранения созданной должности
func addPosition(ctx context.Context, b *bot.Bot, update *models.Update) {
	var answerTxt string
	position := pstn.Position{PositionName: update.Message.Text}
	if err := position.Insert(ctx); err != nil {
		answerTxt = "❌ В процессе создания должности произошла ошибка: " + err.Error()
	} else {
		answerTxt = "✅ Должность успешно добавлена"
		state.Del(context.TODO(), "user_state", strconv.FormatInt(update.Message.From.ID, 10))
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      highlightTxt(answerTxt),
		ParseMode: models.ParseModeHTML,
	})
	StartHandler(ctx, b, update)
}

// positionShowHandler функция для вывода существующих должностей
func positionShowHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	positions, _ := pstn.SelectAllStr(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Список должностей:" + "\n" + strings.Join(positions, "\n")),
		ReplyMarkup: inline.New(b).Button("Назад", []byte(""), positionSettingHandler),
		ParseMode:   models.ParseModeHTML,
	})
}
