package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/slider"
	emp "go_telegram_bot/src/database/employee"
	pstn "go_telegram_bot/src/database/position"
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
		"employeePhotoAdd":       "Приложите фото",
		"empPositionAdd":         "Выберите должность сотрудника",
	}
)

// EmpSettingHandler функция для формирования настроек сотрудников
func empSettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	kb := inline.New(b).
		Row().
		Button("Добавить сотрудника", []byte(""), empAddStart).
		Row().
		Button("Удалить сотрудника", []byte(""), empDeleteStart).
		Row().
		Button("Вывести сотрудников", []byte(""), empShowHandler).
		Row().
		Button("⬅️ Назад", []byte(""), BackSettingHandler)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// empAddStart функция для начала процесса создания сотрудника
func empAddStart(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	empAddAttr(ctx, b, mes, "employeeFirstNameAdd")
}

// empAddAttr функция для смены промежуточного состояния процесса создания сотрудника
func empAddAttr(ctx context.Context, b *bot.Bot, mes *models.Message, empState string) {
	state.Set(ctx, "user_state", strconv.FormatInt(mes.Chat.ID, 10), empState)
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt(empText[empState]),
		ParseMode: models.ParseModeHTML,
	})
}

// empPositionAdd функция формирования клавиатуры выбора должности сотрудника со сменой состония
func empPositionAdd(ctx context.Context, b *bot.Bot, mes *models.Message, empState string) {
	state.Set(ctx, "user_state", strconv.Itoa(mes.ID), empState)
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
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
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
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
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt(res),
		ParseMode: models.ParseModeHTML,
	})
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
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
func empDeleteStart(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите сотрудника на удаления"),
		ReplyMarkup: empKB(ctx, b, empDelete).Row().Button("Назад", []byte(""), empSettingHandler),
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
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt(res),
		ParseMode: models.ParseModeHTML,
	})

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: menuKeyboard(b),
		ParseMode:   models.ParseModeHTML,
	})
}

// empShowHandler функция для вывода существующих сотрудников
func empShowHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	slides, err := getEmpSlides(ctx)
	if err != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      mes.Chat.ID,
			Text:        highlightTxt("Ошибка в процессе формирования списка сотрудников!"),
			ReplyMarkup: inline.New(b).Button("Назад", []byte(""), empSettingHandler),
			ParseMode:   models.ParseModeHTML,
		})
	}
	opts := []slider.Option{
		//slider_cust.OnSelect("Select", true, sliderOnSelect),
		slider.OnCancel("Назад", true, sliderOnCancel),
	}
	sl := slider.New(slides, opts...)
	m, err := sl.Show(ctx, b, mes.Chat.ID)
	fmt.Println(m, err)
}

// getEmpSlides функция для формирования клавиатуры слайдов
func getEmpSlides(ctx context.Context) ([]slider.Slide, error) {
	var slides []slider.Slide
	employees, err := emp.SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, employee := range employees {
		slides = append(slides, slider.Slide{
			Photo:    string(employee.Photo),
			IsUpload: true,
			Text:     employee.MiddleName + " " + employee.FirstName + " " + employee.LastName,
		})
	}
	return slides, nil
}

// sliderOnCancel функция возврата в настройки
func sliderOnCancel(ctx context.Context, b *bot.Bot, mes *models.Message) {
	empSettingHandler(ctx, b, mes, []byte(""))
}
