package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
	sttng "go_telegram_bot/database/setting"
	"time"
)

func CalendarHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := datepicker.New(b, onDatepickerSimpleSelect, datepicker.Language("ru"))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select any date",
		ReplyMarkup: kb,
	})
}

func getStrSched(startTime time.Time, endTime time.Time) string {
	return startTime.Format("15:04") + "-" + endTime.Format("15:04")
}

func onDatepickerSimpleSelect(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	kbTime := inline.New(b)
	curDayName := time.Now().Format("Mon")
	startTime := date.Add(time.Hour * time.Duration(sttng.StartHourScheduler[curDayName]))
	endTime := date.Add(time.Hour * time.Duration(sttng.EndHourScheduler[curDayName]))
	rowWidthCnt := 0
	for {
		if startTime.After(endTime) {
			break
		}
		nextTime := startTime.Add(time.Minute * time.Duration(60*sttng.SessionTimeHour))
		fmt.Println(nextTime)
		kbTime.Button(getStrSched(startTime, nextTime), []byte("1-1"), onInlineKeyboardSelect)
		rowWidthCnt += 1
		if rowWidthCnt == sttng.TimeKeyboarWidth {
			kbTime.Row()
			rowWidthCnt = 0
		}
		startTime = nextTime
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        "Select the variant",
		ReplyMarkup: kbTime,
	})
}

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You selected: " + string(data),
	})
}
