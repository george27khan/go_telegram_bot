package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
	schdlr "go_telegram_bot/database/schedule"
	sttng "go_telegram_bot/database/setting"
	"time"
)

const datetimeFormat string = "02.01.2006 15:04"

func CalendarHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {

	//excludeDays := []time.Time{
	//	makeTime(2020, 1, 10),
	//	makeTime(2020, 1, 13),
	//	makeTime(2019, 12, 27),
	//	makeTime(2019, 12, 28),
	//	makeTime(2019, 12, 29),
	//}
	opts := []datepicker.Option{
		datepicker.StartFromSunday(),
		datepicker.CurrentDate(time.Now()),
		datepicker.From(time.Now()),
		datepicker.To(time.Now().AddDate(0, 0, sttng.DaysInSchedule)),
		//datepicker.OnCancel(onDatepickerCustomCancel),
		datepicker.Language(lang),
		//datepicker.Dates(datepicker.DateModeExclude, excludeDays),
	}

	kb := datepicker.New(b, TimeHandler, opts...)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите дату записи:"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

func getStrSched(startTime time.Time, endTime time.Time) string {
	return startTime.Format("15:04") + "-" + endTime.Format("15:04")
}

func TimeHandler(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	kbTime := inline.New(b)
	curDayName := time.Now().Format("Mon")
	startTime := date.Add(time.Hour * time.Duration(sttng.StartHourSchedule[curDayName]))
	endTime := date.Add(time.Hour * time.Duration(sttng.EndHourSchedule[curDayName]))
	rowWidthCnt := 0
	for {
		nextTime := startTime.Add(time.Minute * time.Duration(60*sttng.SessionTimeHour))
		kbTime.Button(getStrSched(startTime, nextTime), []byte(startTime.Format(datetimeFormat)), TimeAnswerHandler)
		rowWidthCnt += 1
		if rowWidthCnt == sttng.TimeKeyboarWidth {
			kbTime.Row()
			rowWidthCnt = 0
		}
		startTime = nextTime
		if startTime.After(endTime) || startTime.Equal(endTime) {
			break
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберете время записи:"),
		ReplyMarkup: kbTime,
		ParseMode:   models.ParseModeHTML,
	})
}

func TimeAnswerHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var sendMsg string
	if schedTime, err := time.Parse(datetimeFormat, string(data)); err != nil {
		sendMsg = highlightTxt("В процессе записи произошла ощибка: " + err.Error())
	} else {
		if err := schdlr.InsertSchedule(ctx, mes.Chat.ID, schedTime); err != nil {
			sendMsg = highlightTxt("В процессе записи произошла ощибка: " + err.Error())
		} else {
			sendMsg = highlightTxt("Вы записались на " + string(data))
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      sendMsg,
		ParseMode: models.ParseModeHTML})
}
