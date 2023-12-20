package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	emp "go_telegram_bot/database/employee"
	pstn "go_telegram_bot/database/position"
	"go_telegram_bot/src/state"
	"strconv"
)

var (
	// справочник для хранения текста сообщений состояний создания сотрудника
	empText = map[string]string{
		"employeeFirstNameAdd":   "Введите имя сотрудника",
		"employeeMiddleNameAdd":  "Введите фамилию сотрудника",
		"employeeLastNameAdd":    "Введите отчетсво сотрудника",
		"employeeBirthDateAdd":   "Введите дату рождения сотрудника (дд.мм.гггг)",
		"employeeEmailAdd":       "Введите email сотрудника",
		"employeePhoneNumberAdd": "Введите номер телефона сотрудника",
		"employeeHireDateAdd":    "Введите дату приема сотрудника",
		"empPositionAdd":         "Выберите должность сотрудника",
	}
)

// EmpSettingHandler функция для формирования настроек сотрудников
func EmpSettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	kb := inline.New(b).
		Row().
		Button("Добавить сотрудника", []byte(""), empAddStart).
		Row().
		Button("Удалить сотрудника", []byte(""), empDeleteStart).
		Row().
		Button("Назад", []byte(""), BackSettingHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// empAddStart функция для начала процесса создания сотрудника
func empAddStart(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	empAddAttr(ctx, b, mes, "employeeFirstNameAdd")
}

// empAddAttr функция для смены промежуточного состояния процесса создания сотрудника
func empAddAttr(ctx context.Context, b *bot.Bot, mes *models.Message, empState string) {
	state.Set(ctx, "user_state", strconv.FormatInt(mes.Chat.ID, 10), empState)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt(empText[empState]),
		ParseMode: models.ParseModeHTML,
	})
}

// empPositionAdd функция формирования клавиатуры выбора должности сотрудника со сменой состония
func empPositionAdd(ctx context.Context, b *bot.Bot, mes *models.Message, empState string) {
	state.Set(ctx, "user_state", strconv.Itoa(mes.ID), empState)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt(empText[empState]),
		ReplyMarkup: positionNameKB(ctx, b, ChoosePosition),
		ParseMode:   models.ParseModeHTML,
	})
}

// ChoosePosition функция сохранения выбранной должности сотрудника и сохранение записи
func ChoosePosition(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var (
		res string
	)
	id, _ := strconv.Atoi(string(data))
	if position, err := pstn.Get(ctx, id); err == nil {
		empCash[mes.Chat.ID].Position = position
	} else {
		res = "❌ Ошибка в процессе выбора должности: " + err.Error()
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt(res),
			ParseMode: models.ParseModeHTML,
		})
	}
	if err := empCash[mes.Chat.ID].Insert(ctx); err != nil {
		res = "❌ Ошибка в процессе сохранения сотрудника: " + err.Error()
	} else {
		res = "✅ Сотрудник успешно добавлен!"
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

// empKB функция формирующая клавиатуру со списком сотрудников
func empKB(ctx context.Context, b *bot.Bot, action func(context.Context, *bot.Bot, *models.Message, []byte)) *inline.Keyboard {
	kb := inline.New(b)
	employees, err := emp.SelectAll(ctx)
	if err != nil {
		fmt.Println("empKB error:", err)
	}
	for _, empl := range employees {
		kb = kb.Row().Button(empl.MiddleName+" "+empl.FirstName+" "+empl.LastName, []byte(strconv.Itoa(empl.Id)), action)
	}
	return kb
}

// empDeleteStart функция для вывода списка сотрудников для удаления
func empDeleteStart(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите должность на удаления"),
		ReplyMarkup: empKB(ctx, b, empDelete),
		ParseMode:   models.ParseModeHTML,
	})
}

// empDeleteStart функция удаления из базы выбранного сотрудника
func empDelete(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	var (
		res string
	)
	id, err := strconv.Atoi(string(data))
	if err != nil {
		res = "❌ Ошибка в процессе удаления сотрудника: " + err.Error()
	}
	if err := emp.DeleteById(ctx, id); err != nil {
		res = "❌ Ошибка в процессе удаления сотрудника: " + err.Error()
	} else {
		res = "✅ Сотрудник удален"
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
