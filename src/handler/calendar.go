package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
	schdlr "go_telegram_bot/database/schedule"
	sttng "go_telegram_bot/database/setting"
	"go_telegram_bot/src/slider_cust"
	"strconv"
	"time"
)

const datetimeFormat string = "02.01.2006 15:04"

var (
	schedTimeCash = map[int64]*schdlr.Schedule{}
)

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
		kbTime.Button(getStrSched(startTime, nextTime), []byte(startTime.Format(datetimeFormat)), schedEmpHandler)
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

func schedEmpHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var (
		slides []slider_cust.Slide
	)
	schedTime, errParse := time.Parse(datetimeFormat, string(data))
	if errParse != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      mes.Chat.ID,
			Text:        highlightTxt("В процессе преобразования даты произошла ощибка: " + errParse.Error()),
			ReplyMarkup: inline.New(b).Button("Назад", []byte(""), empSettingHandler),
			ParseMode:   models.ParseModeHTML,
		})
		return
	}

	empSlice, errSched := schdlr.GetFreeEmpVisitDt(ctx, schedTime)
	if errSched != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      mes.Chat.ID,
			Text:        highlightTxt("В процессе получения свободных сотрудников произошла ощибка: " + errSched.Error()),
			ReplyMarkup: inline.New(b).Button("Назад", []byte(""), empSettingHandler),
			ParseMode:   models.ParseModeHTML,
		})
	}

	for _, employee := range empSlice {
		slides = append(slides, slider_cust.Slide{
			Photo:    string(employee.Photo),
			IsUpload: true,
			Text:     employee.MiddleName + " " + employee.FirstName + " " + employee.LastName,
			Data:     []byte(strconv.Itoa(employee.Id)),
		})
	}

	opts := []slider_cust.Option{
		slider_cust.OnSelect("Выбрать", true, sliderOnSelect),
		slider_cust.OnCancel("Назад", true, sliderOnCancel1),
	}
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt("Выберите сотрудника"),
		ParseMode: models.ParseModeHTML,
	})

	sl := slider_cust.New(slides, opts...)
	_, _ = sl.Show(ctx, b, mes.Chat.ID)
	schedTimeCash[mes.Chat.ID] = &schdlr.Schedule{IdUser: mes.Chat.ID, VisitDt: schedTime}
}

func sliderOnSelect(ctx context.Context, b *bot.Bot, mes *models.Message, item int, data []byte) {
	if idEmp, err := strconv.ParseInt(string(data), 10, 64); err != nil {
		fmt.Println("Error")
	} else {
		fmt.Println(schedTimeCash[mes.Chat.ID].VisitDt, idEmp)
		schedTimeCash[mes.Chat.ID].IdEmployee = idEmp
	}
	if err := schedTimeCash[mes.Chat.ID].Insert(ctx); err != nil {
		fmt.Println("Error")
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "Select " + string(data),
	})
}

// sliderOnCancel1 функция возврата в настройки
func sliderOnCancel1(ctx context.Context, b *bot.Bot, mes *models.Message) {
	empSettingHandler(ctx, b, mes, []byte(""))
}
