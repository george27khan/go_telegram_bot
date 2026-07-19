package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	usr "go_telegram_bot/internal/infrastructure/repository/postgres/client"
	fill "go_telegram_bot/internal/infrastructure/repository/postgres/fill_table"
	"go_telegram_bot/src/Petrovich"
	emp "go_telegram_bot/src/database/employee"
	"go_telegram_bot/src/state"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"time"
)

// формат даты
const dateFormat string = "02.01.2006"

// шаблон url для бота
const fileDownloadURL string = "https://api.telegram.org/file/bot%s/%s" //https://api.telegram.org/file/bot<token>/<file_path>

var (
	lang       string // язык для чата
	empCash    = map[int64]*emp.Employee{}
	declension *Petrovich.Rules
)

// highlightTxt функция выделения текста сообщения телеграмм
func highlightTxt(str string) string {
	return "<b>" + str + "</b>"
}

// menuKeyboard функция формирования клавиатуры меню
func menuKeyboard(b *bot.Bot) *inline.Keyboard {
	return inline.New(b).
		Row().
		Button("📆 Запись на прием", []byte(""), CalendarHandler).
		Row().
		Button("📝 Записи", []byte(""), entryHandler).
		Row().
		Button("⚙️ Настройки", []byte(""), settingHandler).
		Row().
		Button("❌ Выход", []byte(""), cancelHandler)
}

// StartHandler функция вывода начального меню
func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.LanguageCode == "ru" {
		lang = "ru"
	} else {
		lang = "en"
	}
	if isExist, err := usr.IsExists(ctx, update.Message.From.Username); err != nil {
		fmt.Println("Ошибка при поиска бзера ", err.Error())
	} else {
		if !isExist {
			if err := usr.Insert(ctx, update.Message.From.ID, update.Message.From.Username, update.Message.From.FirstName, update.Message.From.LastName, ""); err != nil {
				fmt.Println("Ошибка при регистрации юзера ", err.Error())
			}
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: menuKeyboard(b),
		ParseMode:   models.ParseModeHTML,
	})
}

// settingHandler функция вывода меню настроек
func settingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	kb := inline.New(b).
		Row().
		Button("Должность", []byte(""), positionSettingHandler).
		Row().
		Button("Сотрудник", []byte(""), empSettingHandler).
		Row().
		Button("Заполнить данные", []byte(""), empInitHandler).
		Row().
		Button("⬅️ Назад", []byte(""), BackStartHandler)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      mes.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// empInitHandler функция для заполнения тестовыми данными
func empInitHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	if err := fill.AddPositions(ctx); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	if err := fill.AddEmployees(ctx); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    mes.Chat.ID,
			Text:      highlightTxt("Ошибка в процессе заполнения таблиц " + err.Error()),
			ParseMode: models.ParseModeHTML,
		})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    mes.Chat.ID,
		Text:      highlightTxt("Таблицы успешно заполнены"),
		ParseMode: models.ParseModeHTML,
	})
}

func downloadFile(URL string) ([]byte, error) {
	//Get the response bytes from the url
	var bytes []byte
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		bytes, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

// DefaultHandler процедура для обработки произвольного сообщения пользователя по текущему состоянию
func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	userIdStr := strconv.FormatInt(update.Message.From.ID, 10)
	userId := update.Message.From.ID
	res := state.Get(ctx, "user_state", userIdStr)
	fmt.Println(res)
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
			empCash[userId].BirthDate = date
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
		empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
	} else if res == "employeePhotoAdd" {
		var (
			fileId string
			url    string
		)
		if update.Message.Document != nil && update.Message.Document.MimeType == "image/jpeg" {
			fileId = update.Message.Document.FileID
		} else if update.Message.Photo != nil {
			fileId = update.Message.Photo[len(update.Message.Photo)-1].FileID
		}
		if fileId != "" {
			if file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileId}); err == nil {
				url = fmt.Sprintf(fileDownloadURL, "5844620699:AAGbEPIFWKxTDr0jR_A77Rba95jtZBSQlGM", file.FilePath)
				if photo, err := downloadFile(url); err == nil {
					empCash[userId].Photo = photo
					empAddAttr(ctx, b, update.Message, "employeeHireDateAdd")
				} else {
					empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
				}
			} else {
				empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
			}
		} else {
			empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
		}
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
func cancelHandler(_ context.Context, _ *bot.Bot, _ *models.Message, _ []byte) {}

// BackStartHandler функция возврата в основное меню
func BackStartHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "/start",
	})
}

// BackSettingHandler функция возврата в меню настроек
func BackSettingHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	settingHandler(ctx, b, mes, data)
}

func init() {
	declension, _ = Petrovich.LoadRules("./src/Petrovich/rules.json")
}
