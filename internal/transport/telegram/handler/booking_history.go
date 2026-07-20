package handler

// import (
// 	"context"
// 	"github.com/go-telegram/bot"
// 	"github.com/go-telegram/bot/models"
// 	"github.com/go-telegram/ui/keyboard/inline"
// 	"go_telegram_bot/internal/domain/entity"
// 	"go_telegram_bot/src/Petrovich"
// 	"log/slog"
// 	"time"
// )

// type SchedulerEmployeeHandler struct {
// 	//CalendarHandler CalendarHandler
// }

// type BookingHistiryUseCase interface {
// 	GetUserBooking(ctx context.Context, id entity.ClientID) (dateFrom time.Time, dateTo time.Time, err error)
// 	GetUserBookingHist(ctx context.Context, id entity.ClientID) (dateFrom time.Time, dateTo time.Time, err error)
// }

// type BookingHistiryHandler struct {
// 	BookingHistiryUseCase BookingHistiryUseCase
// 	logger                *slog.Logger
// }

// func NewBookingHistiryHandler(bookingHistiryUseCase BookingHistiryUseCase, logger *slog.Logger) BookingHistiryHandler {
// 	return BookingHistiryHandler{
// 		BookingHistiryUseCase: bookingHistiryUseCase,
// 		logger:                logger,
// 	}
// }

// func (bh *BookingHistiryHandler) Get(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
// 	kb := inline.New(b).
// 		Row().
// 		Button("Текущая запись", []byte(""), bh.actualEntryHandler).
// 		Row().
// 		Button("История записей", []byte(""), bh.histEntryHandler).
// 		Row().
// 		Button("⬅️ Назад", []byte(""), BackStartHandler)

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      mes.Chat.ID,
// 		Text:        highlightTxt("Выберите действие"),
// 		ReplyMarkup: kb,
// 		ParseMode:   models.ParseModeHTML,
// 	})
// }

// // actualEntryHandler функция вывода актуальной записи
// func (bh *BookingHistiryHandler) actualEntryHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
// 	schedule, err := sched.GetByUser(ctx, mes.Chat.ID)
// 	if err != nil {
// 		b.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID:    mes.Chat.ID,
// 			Text:      highlightTxt("Ошибка в процессе поиска записи " + err.Error()),
// 			ParseMode: models.ParseModeHTML,
// 		})
// 		return
// 	}
// 	FIO, err := emp.GetFIO(ctx, schedule.IdEmployee)
// 	if err != nil {
// 		b.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID:    mes.Chat.ID,
// 			Text:      highlightTxt("Ошибка в процессе поиска сотрудника " + err.Error()),
// 			ParseMode: models.ParseModeHTML,
// 		})
// 		return
// 	}

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      mes.Chat.ID,
// 		Text:        highlightTxt("Вы записаны к " + declension.InfFio(FIO, Petrovich.Dative, false) + " на " + schedule.VisitDt.Format(datetimeFormat)),
// 		ReplyMarkup: inline.New(b).Button("⬅️ Назад", []byte(""), entryHandler),
// 		ParseMode:   models.ParseModeHTML,
// 	})
// }

// // histEntryHandler функция вывод истории записей
// func (bh *BookingHistiryHandler) histEntryHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
// 	var (
// 		res string
// 	)
// 	schedSlice, err := sched.GetAllByUser(ctx, mes.Chat.ID)
// 	if err != nil {
// 		b.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID:    mes.Chat.ID,
// 			Text:      highlightTxt("Ошибка в процессе поиска записей " + err.Error()),
// 			ParseMode: models.ParseModeHTML,
// 		})
// 		return
// 	}
// 	for _, schedule := range schedSlice {
// 		FIO, err := emp.GetFIO(ctx, schedule.IdEmployee)
// 		if err != nil {
// 			b.SendMessage(ctx, &bot.SendMessageParams{
// 				ChatID:    mes.Chat.ID,
// 				Text:      highlightTxt("Ошибка в процессе поиска сотрудника " + err.Error()),
// 				ParseMode: models.ParseModeHTML,
// 			})
// 			return
// 		} else {
// 			res = res + schedule.VisitDt.Format(datetimeFormat) + " " + FIO + "\n"
// 		}
// 	}

// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID:      mes.Chat.ID,
// 		Text:        highlightTxt(res),
// 		ReplyMarkup: inline.New(b).Button("⬅️ Назад", []byte(""), entryHandler),
// 		ParseMode:   models.ParseModeHTML,
// 	})
// }
