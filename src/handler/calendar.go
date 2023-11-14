package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"
	"github.com/go-telegram/ui/keyboard/inline"
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

func onDatepickerSimpleSelect(ctx context.Context, b *bot.Bot, mes *models.Message, date time.Time) {
	g_hour_start_d := map[string]int{"Mon": 9, "Tue": 9, "Wed": 9, "Thu": 9, "Fri": 9, "Sat": 9, "Sun": 9}
	g_hour_end_d := map[string]int{"Mon": 18, "Tue": 18, "Wed": 18, "Thu": 18, "Fri": 18, "Sat": 18, "Sun": 18}
	curDayName := time.Now().Format("Mon")
	curDayName := time.Now().Format("Mon")
	startHour = g_hour_start_d[curDayName]
	endHour = g_hour_end_d[curDayName]

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "You select " + date.Format("2006-01-02"),
	})

	kb := inline.New(b).
		Row().
		Button("", []byte("1-1"), onInlineKeyboardSelect).
		Button("Row 1, Btn 2", []byte("1-2"), onInlineKeyboardSelect).
		Row().
		Button("Row 2, Btn 1", []byte("2-1"), onInlineKeyboardSelect).
		Button("Row 2, Btn 2", []byte("2-2"), onInlineKeyboardSelect).
		Button("Row 2, Btn 3", []byte("2-3"), onInlineKeyboardSelect).
		Row().
		Button("Row 3, Btn 1", []byte("3-1"), onInlineKeyboardSelect).
		Button("Row 3, Btn 2", []byte("3-2"), onInlineKeyboardSelect).
		Button("Row 3, Btn 3", []byte("3-3"), onInlineKeyboardSelect).
		Button("Row 3, Btn 4", []byte("3-4"), onInlineKeyboardSelect).
		Row().
		Button("Cancel", []byte("cancel"), onInlineKeyboardSelect)


	if call.data in ("empty_day", "week_day") or "day_" not in call.data:
	await call.answer()
	return
	button_datetime = dt.datetime.strptime(call.data[4:], g_date_format)
	await Booking.choose_time.set()  # встаем в состояние выбора дня
	await call.message.answer(
		"Выберите время бронирования:",
		reply_markup=get_time_keyboard(button_datetime),
)
	await call.answer()



	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        "Select the variant",
		ReplyMarkup: kb,
	})
}
