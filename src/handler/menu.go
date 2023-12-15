package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	pstn "go_telegram_bot/database/position"
	redis "go_telegram_bot/database/redis"
	usr "go_telegram_bot/database/user"
	"strconv"
)

var (
	lang         string            // язык для чата
	user_action  string            // текущее действие пользователя для обработки
	menu_state   map[string]string = map[string]string{"PositionHandler": "position_name"}
	menuKeyboard *inline.Keyboard
)

func getMenuKeyboard(b *bot.Bot) *inline.Keyboard {
	return inline.New(b).
		Row().
		Button("📆 Запись на прием", []byte(""), CalendarHandler).
		Row().
		Button("⚙️ Настройки", []byte(""), SettingHandler).
		Row().
		Button("Cancel", []byte(""), CancelHandler)
}
func highlightTxt(str string) string {
	return "<b>" + str + "</b>"
}

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
		ReplyMarkup: getMenuKeyboard(b),
		ParseMode:   models.ParseModeHTML,
	})
}

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var (
		answerTxt string
	)
	userId := strconv.FormatInt(update.Message.From.ID, 10)
	fmt.Println("DefaultHandler")
	rdb := redis.Connect()
	res := rdb.HGet(context.TODO(), "user_state", userId)
	if res.Val() == "PositionNameHandler" {
		if err := pstn.InsertPosition(ctx, update.Message.Text); err != nil {
			answerTxt = "❌ В процессе создания должности произошла ошибка: " + err.Error()
		} else {
			answerTxt = "✅ Должность успешно добавлена"
			_ = rdb.HDel(context.TODO(), "user_state", userId)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      highlightTxt(answerTxt),
			ParseMode: models.ParseModeHTML,
		})
		StartHandler(ctx, b, update)
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      highlightTxt("Выберите дейтсиве из меню"),
			ParseMode: models.ParseModeHTML,
		})
	}

}

func CancelHandler(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {}
