package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	pstn "go_telegram_bot/database/position"
	redis "go_telegram_bot/database/redis"
	"strconv"
	"strings"
)

func SettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b).
		Row().
		Button("Должность", []byte("/schedule"), PositionHandler).
		Row().
		Button("Cотрудник", []byte("2-1"), CalendarHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

func PositionHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b).
		Row().
		Button("Добавить должность", []byte(""), PositionAddHandler).
		Row().
		Button("Удалить должность", []byte(""), PositionDelHandler).
		Row().
		Button("Вывести должности", []byte(""), PositionShowHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

func PositionDelHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b)
	positionNames, err := pstn.SelectAllPositionMap(ctx)
	if err != nil {
		fmt.Println("PositionDelHandler error:", err)
	}
	for id, name := range positionNames {
		kb = kb.Row().Button(name, []byte(strconv.Itoa(id)), DelPosition)

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите должность на удаления"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

func DelPosition(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var (
		res string
	)
	id, err := strconv.Atoi(string(data))
	if err != nil {
		res = "❌ Ошибка в процессе удаления позиции: " + err.Error()
	}
	if err := pstn.DeletePositionById(ctx, id); err != nil {
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
		ReplyMarkup: getMenuKeyboard(b),
		ParseMode:   models.ParseModeHTML,
	})
}

func PositionAddHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	rdb := redis.Connect()
	defer rdb.Close()
	_ = rdb.HSet(ctx, "user_state", map[string]interface{}{strconv.FormatInt(mes.Chat.ID, 10): "PositionNameHandler"})
	positions, _ := pstn.SelectAllPosition(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt("Введите новое название должности НЕ из списка:" + "\n" + strings.Join(positions, "\n")),
		ParseMode: models.ParseModeHTML,
	})
}

func PositionShowHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	rdb := redis.Connect()
	defer rdb.Close()
	_ = rdb.HSet(ctx, "user_state", map[string]interface{}{strconv.FormatInt(mes.Chat.ID, 10): "PositionNameHandler"})
	positions, _ := pstn.SelectAllPosition(ctx)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt("Список должностей:" + "\n" + strings.Join(positions, "\n")),
		ParseMode: models.ParseModeHTML,
	})
}
