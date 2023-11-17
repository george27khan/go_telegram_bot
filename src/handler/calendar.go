package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
	"time"
)

func CalendarHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := datepicker.New(b, onDatepickerSimpleSelect, datepicker.Language("ru"))
	fmt.Print(update.Message.Chat.ID)
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
	keyboarWidth := 3
	kbTime := inline.New(b)
	hourStartM := map[string]int{"Mon": 9, "Tue": 9, "Wed": 9, "Thu": 9, "Fri": 9, "Sat": 9, "Sun": 9}
	hourEndM := map[string]int{"Mon": 18, "Tue": 18, "Wed": 18, "Thu": 18, "Fri": 18, "Sat": 18, "Sun": 18}
	timeStep := 15
	curDayName := time.Now().Format("Mon")

	endHour := hourEndM[curDayName]
	startTime := date.Add(time.Hour * time.Duration(hourStartM[curDayName]))
	rowWidthCnt := 0
	for {
		if startTime.Hour() >= endHour {
			break
		}
		nextTime := startTime.Add(time.Minute * time.Duration(timeStep))

		kbTime.Button(getStrSched(startTime, nextTime), []byte("1-1"), onInlineKeyboardSelect)
		rowWidthCnt += 1
		if rowWidthCnt == keyboarWidth {
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
