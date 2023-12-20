package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	emp "go_telegram_bot/database/employee"
	usr "go_telegram_bot/database/user"
	"go_telegram_bot/src/state"
	"net/mail"
	"strconv"
	"time"
)

const dateFormat string = "02.01.2006"

var (
	lang    string // язык для чата
	empCash = map[int64]*emp.Employee{}
)

// menuKeyboard функция формирования клавиатуры меню
func menuKeyboard(b *bot.Bot) *inline.Keyboard {
	return inline.New(b).
		Row().
		Button("📆 Запись на прием", []byte(""), CalendarHandler).
		Row().
		Button("⚙️ Настройки", []byte(""), settingHandler).
		Row().
		Button("Cancel", []byte(""), cancelHandler)
}

// highlightTxt функция выделения текста сообщения телеграмм
func highlightTxt(str string) string {
	return "<b>" + str + "</b>"
}

// StartHandler функция вывода начального меню
func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.LanguageCode == "ru" {
		lang = "ru"
	} else {
		lang = "en"
	}
	if err := usr.InsertUser(ctx, update.Message.From.ID, update.Message.From.Username, update.Message.From.FirstName, update.Message.From.LastName, ""); err != nil {
		fmt.Println(err)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: menuKeyboard(b),
		ParseMode:   models.ParseModeHTML,
	})
}

// settingHandler функция вывода меню настроек
func settingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b).
		Row().
		Button("Должность", []byte(""), positionHandler).
		Row().
		Button("Сотрудник", []byte(""), EmpSettingHandler).
		Row().
		Button("Назад", []byte(""), BackStartHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// DefaultHandler процедура для обработки произвольного сообщения пользователя по текущему состоянию
func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userIdStr := strconv.FormatInt(update.Message.From.ID, 10)
	userId := update.Message.From.ID
	res := state.Get(ctx, "user_state", userIdStr)
	if res == "PositionAddHandler" {
		addPosition(ctx, b, update)
	} else if res == "employeeFirstNameAdd" {
		employee := new(emp.Employee)
		employee.FirstName = update.Message.Text
		empCash[userId] = employee
		empAddAttr(ctx, b, update.Message, "employeeMiddleNameAdd")
	} else if res == "employeeMiddleNameAdd" {
		empCash[userId].MiddleName = update.Message.Text
		empAddAttr(ctx, b, update.Message, "employeeLastNameAdd")
	} else if res == "employeeLastNameAdd" {
		empCash[userId].LastName = update.Message.Text
		empAddAttr(ctx, b, update.Message, "employeeBirthDateAdd")
	} else if res == "employeeBirthDateAdd" {
		if date, err := time.Parse(dateFormat, update.Message.Text); err == nil {
			empCash[userId].BithDate = date
			empAddAttr(ctx, b, update.Message, "employeeEmailAdd")
		} else {
			empAddAttr(ctx, b, update.Message, "employeeBirthDateAdd")
		}
	} else if res == "employeeEmailAdd" {
		if _, err := mail.ParseAddress(update.Message.Text); err == nil {
			empCash[userId].Email = update.Message.Text
			empAddAttr(ctx, b, update.Message, "employeePhoneNumberAdd")
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      highlightTxt("Введенный адрес не прошел валидацию."),
				ParseMode: models.ParseModeHTML,
			})
			empAddAttr(ctx, b, update.Message, "employeeEmailAdd")
		}
	} else if res == "employeePhoneNumberAdd" {
		empCash[userId].PhoneNumber = update.Message.Text
		empAddAttr(ctx, b, update.Message, "employeeHireDateAdd")

	} else if res == "employeeHireDateAdd" {
		if date, err := time.Parse(dateFormat, update.Message.Text); err == nil {
			empCash[userId].HireDate = date
			empPositionAdd(ctx, b, update.Message, "empPositionAdd")
		} else {
			fmt.Println(err)
			empAddAttr(ctx, b, update.Message, "employeeHireDateAdd")
		}
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      highlightTxt("Выберите дейтсиве из меню"),
			ParseMode: models.ParseModeHTML,
		})
	}

}

// cancelHandler пустая функция для выхода из меню
func cancelHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {}

// BackStartHandler функция возврата в основное меню
func BackStartHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "/start",
	})
}

// BackSettingHandler функция возврата в меню настроек
func BackSettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	settingHandler(ctx, b, mes, data)
}
