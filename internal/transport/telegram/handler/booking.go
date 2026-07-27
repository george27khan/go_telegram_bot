package handler

import (
	"context"
	"fmt"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/pkg/slider_cust"
	"log"
	"log/slog"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"

	//"go_telegram_bot/src/Petrovich"
	"time"
)

const datetimeFormat string = "02.01.2006 15:04"

var (
	schedTimeCash = map[int64]*entity.Schedule{}
)

type BookingUseCase interface {
	GetCalendarDays(ctx context.Context) (dateFrom time.Time, dateTo time.Time, err error)
	GetDayTimes(ctx context.Context, date time.Time) (timeFrom time.Time, timeTo time.Time, err error)
	GetSessionTimeMinute(ctx context.Context) (time.Duration, error)
	IsFreeTime(ctx context.Context, visitTime time.Time) (ok bool, err error)
	GetFreeEmployee(ctx context.Context, visitTime time.Time) ([]entity.Employee, error)
	SaveSchedule(ctx context.Context, schedule entity.Schedule) error
}

type BookingHandler struct {
	BookingUseCase BookingUseCase
	schedule       entity.Schedule
	logger         *slog.Logger
}

func NewBookingHandler(bookingUseCase BookingUseCase, logger *slog.Logger) BookingHandler {
	return BookingHandler{
		BookingUseCase: bookingUseCase,
		logger:         logger,
	}
}

func getStrSched(startTime time.Time, endTime time.Time) string {
	return startTime.Format("15:04") + "-" + endTime.Format("15:04")
}

// GetCalendarn отрисовка в меню календаря для записи
func (bh *BookingHandler) GetCalendar(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	dateFrom, dateTo, err := bh.BookingUseCase.GetCalendarDays(ctx)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	opts := []datepicker.Option{
		datepicker.CurrentDate(dateFrom),
		datepicker.From(dateFrom),
		datepicker.To(dateTo),
		//datepicker.OnCancel(onDatepickerCustomCancel),
		datepicker.Language(mes.From.LanguageCode),
		//datepicker.Dates(datepicker.DateModeExclude, excludeDays),
	}

	kb := datepicker.New(b, bh.getTime, opts...)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите дату записи:"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// getTime отрисовка меню для выбора времени записи по определенному дню
func (bh *BookingHandler) getTime(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	kbTime := inline.New(b)
	timeFrom, timeTo, err := bh.BookingUseCase.GetDayTimes(ctx, date)
	log.Println("timeFrom, timeTo",timeFrom, timeTo)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	sessionTime, err := bh.BookingUseCase.GetSessionTimeMinute(ctx)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	rowWidthCnt := 0
	for {
		nextTime := timeFrom.Add(sessionTime)
		ok, err := bh.BookingUseCase.IsFreeTime(ctx, timeFrom)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    mes.Chat.ID,
				Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
				ParseMode: models.ParseModeHTML,
			})
			return
		}
		if ok {
			kbTime.Button(getStrSched(timeFrom, nextTime), []byte(timeFrom.Format(datetimeFormat)), bh.getEmployee)
		} else {
			kbTime.Button("-", []byte(""), bh.GetCalendar)
		}
		rowWidthCnt += 1
		if rowWidthCnt == 3 { //sttng.TimeKeyboarWidth
			kbTime.Row()
			rowWidthCnt = 0
		}

		timeFrom = nextTime
		if timeFrom.After(timeTo) || timeFrom.Equal(timeTo) {
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


// getEmployee отрисовка меню для выбора сотрудника для записи
func (bh *BookingHandler) getEmployee(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var (
		slides []slider_cust.Slide
	)
	visitTime, err := time.Parse(datetimeFormat, string(data))
	bh.schedule.VisitDt = visitTime //запоминаем дату и время бронирования
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   highlightTxt("В процессе преобразования даты произошла ощибка: " + err.Error()),
			//ReplyMarkup: inline.New(b).Button("Назад", []byte(""), bh.GetCalendar),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	empList, err := bh.BookingUseCase.GetFreeEmployee(ctx, visitTime)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("В процессе получения свободных сотрудников произошла ошибка: " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
	}

	for _, employee := range empList {
		slides = append(slides, slider_cust.Slide{
			Photo:    string(employee.Photo),
			IsUpload: true,
			Text:     employee.MiddleName + " " + employee.FirstName + " " + employee.LastName,
			Data:     []byte(strconv.Itoa(int(employee.Id))),
		})
	}

	opts := []slider_cust.Option{
		slider_cust.OnSelect("Выбрать", true, bh.saveBooking),
		//slider_cust.OnCancel("Назад", true, ),
	}
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt("Выберите сотрудника"),
		ParseMode: models.ParseModeHTML,
	})

	sl := slider_cust.New(slides, opts...)
	msg, err := sl.Show(ctx, b, mes.Chat.ID)
	fmt.Println(err)
	fmt.Println(msg.Text)
	//schedTimeCash[mes.Chat.ID] = &schdlr.Schedule{IdUser: mes.Chat.ID, VisitDt: schedTime}
}


// saveBooking бронирование записи
func (bh *BookingHandler) saveBooking(ctx context.Context, b *bot.Bot, mes *models.Message, item int, data []byte) {
	idEmp, err := strconv.Atoi(string(data))
	fmt.Println("item", item)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "В процессе получения ИД сотрудника произошла ошибка " + err.Error(),
		})
		return
	}
	bh.schedule.EmployeeID = entity.EmployeeID(idEmp)
	bh.schedule.ClientID = entity.ClientID(mes.Chat.ID)
	if err := bh.BookingUseCase.SaveSchedule(ctx, bh.schedule); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "В процессе сохранения записи произошла ошибка " + err.Error(),
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		//Text:      highlightTxt("Вы записались к " + declension.InfFio(FIO, Petrovich.Dative, false) + " на " + schedTimeCash[mes.Chat.ID].VisitDt.Format(datetimeFormat)),
		Text: highlightTxt("Вы записались на " + bh.schedule.VisitDt.Format(datetimeFormat)),

		ParseMode: models.ParseModeHTML,
	})
}
