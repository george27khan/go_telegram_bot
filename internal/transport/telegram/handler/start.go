package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/transport/telegram/router"
	"log/slog"
)

// шаблон url для бота
const fileDownloadURL string = "https://api.telegram.org/file/bot%s/%s" //https://api.telegram.org/file/bot<token>/<file_path>

//var (
//	empCash    = map[int64]*entity.Employee{}
//	declension *Petrovich.Rules
//)

type StartUseCase interface {
	Start(ctx context.Context, client entity.Client) error
}

type StateUseCase interface {
	Next(current entity.UserState, in entity.UserState) entity.UserState
}

type StartHandler struct {
	StartUseCase StartUseCase
	StateUseCase StateUseCase
	Router       *router.Router
	logger       *slog.Logger
}

func NewStartHandler(startUseCase StartUseCase, stateUseCase StateUseCase, router *router.Router, logger *slog.Logger) *StartHandler {
	return &StartHandler{
		StartUseCase: startUseCase,
		StateUseCase: stateUseCase,
		Router:       router,
		logger:       logger,
	}
}

// highlightTxt функция выделения текста сообщения телеграмм
func highlightTxt(str string) string {
	return "<b>" + str + "</b>"
}

func newClient(user *models.User) entity.Client {
	return entity.Client{
		ID:           entity.ClientID(user.ID),
		UserName:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,
	}
}

// Start функция вывода начального меню
func (sh *StartHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	client := newClient(update.Message.From)
	sh.StartUseCase.Start(ctx, client)

	kb := inline.New(b).
		Row().
		Button("📆 Запись на прием", []byte(""), sh.Router.Route(entity.ActionCalendar)).
		Row().
		Button("📝 Записи", []byte(""), sh.Router.Route(entity.ActionBookHistory)). //entryHandler
		Row().
		Button("⚙️ Настройки", []byte(""), sh.Router.Route(entity.ActionSetting)).
		Row().
		Button("❌ Выход", []byte(""), sh.cancelHandler)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        highlightTxt("Выберите действие"),
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeHTML,
	})
}

// DefaultHandler процедура для обработки произвольного сообщения пользователя по текущему состоянию
//func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
//	userIdStr := strconv.FormatInt(update.Message.From.ID, 10)
//	userId := update.Message.From.ID
//	res := state.Get(ctx, "user_state", userIdStr)
//	fmt.Println(res)
//	if res == "PositionAddHandler" {
//		addPosition(ctx, b, update)
//	} else if res == "employeeFirstNameAdd" {
//		employee := new(emp.Employee)
//		employee.FirstName = update.Message.Text
//		empCash[userId] = employee
//		empAddAttr(ctx, b, update.Message, "employeeMiddleNameAdd")
//	} else if res == "employeeMiddleNameAdd" {
//		empCash[userId].MiddleName = update.Message.Text
//		empAddAttr(ctx, b, update.Message, "employeeLastNameAdd")
//	} else if res == "employeeLastNameAdd" {
//		empCash[userId].LastName = update.Message.Text
//		empAddAttr(ctx, b, update.Message, "employeeBirthDateAdd")
//	} else if res == "employeeBirthDateAdd" {
//		if date, err := time.Parse(dateFormat, update.Message.Text); err == nil {
//			empCash[userId].BirthDate = date
//			empAddAttr(ctx, b, update.Message, "employeeEmailAdd")
//		} else {
//			empAddAttr(ctx, b, update.Message, "employeeBirthDateAdd")
//		}
//	} else if res == "employeeEmailAdd" {
//		if _, err := mail.ParseAddress(update.Message.Text); err == nil {
//			empCash[userId].Email = update.Message.Text
//			empAddAttr(ctx, b, update.Message, "employeePhoneNumberAdd")
//		} else {
//			b.SendMessage(ctx, &bot.SendMessageParams{
//				ChatID:    update.Message.Chat.ID,
//				Text:      highlightTxt("Введенный адрес не прошел валидацию."),
//				ParseMode: models.ParseModeHTML,
//			})
//			empAddAttr(ctx, b, update.Message, "employeeEmailAdd")
//		}
//	} else if res == "employeePhoneNumberAdd" {
//		empCash[userId].PhoneNumber = update.Message.Text
//		empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
//	} else if res == "employeePhotoAdd" {
//		var (
//			fileId string
//			url    string
//		)
//		if update.Message.Document != nil && update.Message.Document.MimeType == "image/jpeg" {
//			fileId = update.Message.Document.FileID
//		} else if update.Message.Photo != nil {
//			fileId = update.Message.Photo[len(update.Message.Photo)-1].FileID
//		}
//		if fileId != "" {
//			if file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: fileId}); err == nil {
//				url = fmt.Sprintf(fileDownloadURL, "5844620699:AAGbEPIFWKxTDr0jR_A77Rba95jtZBSQlGM", file.FilePath)
//				if photo, err := downloadFile(url); err == nil {
//					empCash[userId].Photo = photo
//					empAddAttr(ctx, b, update.Message, "employeeHireDateAdd")
//				} else {
//					empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
//				}
//			} else {
//				empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
//			}
//		} else {
//			empAddAttr(ctx, b, update.Message, "employeePhotoAdd")
//		}
//	} else if res == "employeeHireDateAdd" {
//		if date, err := time.Parse(dateFormat, update.Message.Text); err == nil {
//			empCash[userId].HireDate = date
//			empPositionAdd(ctx, b, update.Message, "empPositionAdd")
//		} else {
//			fmt.Println(err)
//			empAddAttr(ctx, b, update.Message, "employeeHireDateAdd")
//		}
//	} else {
//		b.SendMessage(ctx, &bot.SendMessageParams{
//			ChatID:    update.Message.Chat.ID,
//			Text:      highlightTxt("Выберите дейтсиве из меню"),
//			ParseMode: models.ParseModeHTML,
//		})
//	}
//
//}

// cancelHandler пустая функция для выхода из меню
func (sh *StartHandler) cancelHandler(_ context.Context, _ *bot.Bot, _ *models.Message, _ []byte) {}

// BackStartHandler функция возврата в основное меню
func BackStartHandler(ctx context.Context, b *bot.Bot, mes *models.Message, _ []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   "/start",
	})
}

//func init() {
//	declension, _ = Petrovich.LoadRules("./src/Petrovich/rules.json")
//}
