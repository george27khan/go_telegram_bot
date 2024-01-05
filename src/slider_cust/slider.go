package slider_cust

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnSelectFunc func(ctx context.Context, b *bot.Bot, message *models.Message, item int, date []byte)
type OnCancelFunc func(ctx context.Context, b *bot.Bot, message *models.Message)
type OnErrorFunc func(err error)

type Slide struct {
	Photo    string
	IsUpload bool
	Text     string
	Data     []byte
}

var (
	cmdPrev   = "prev"
	cmdNext   = "next"
	cmdNop    = "nop"
	cmdSelect = "select"
	cmdCancel = "cancel"
)

type Slider struct {
	prefix string
	slides []Slide

	selectButtonText string
	onSelect         OnSelectFunc
	cancelButtonText string
	onCancel         OnCancelFunc
	onError          OnErrorFunc

	deleteOnSelect bool
	deleteOnCancel bool

	current           int
	callbackHandlerID string
}

func New(slides []Slide, opts ...Option) *Slider {
	s := &Slider{
		prefix:  bot.RandomString(16),
		slides:  slides,
		onError: defaultOnError,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func defaultOnError(err error) {
	log.Printf("[TG-UI-SLIDER] [ERROR] %s", err)
}

func (s *Slider) Show(ctx context.Context, b *bot.Bot, chatID any) (*models.Message, error) {
	s.callbackHandlerID = b.RegisterHandler(bot.HandlerTypeCallbackQueryData, s.prefix, bot.MatchTypePrefix, s.callback)

	slide := s.slides[s.current]

	sendParams := &bot.SendPhotoParams{
		ChatID:      chatID,
		Photo:       &models.InputFileString{Data: slide.Photo},
		Caption:     slide.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: s.buildKeyboard(),
	}

	if slide.IsUpload {
		sendParams.Photo = &models.InputFileUpload{
			Filename: "image.png",
			Data:     strings.NewReader(slide.Photo),
		}
	}

	return b.SendPhoto(ctx, sendParams)
}
