package main

import (
	"context"
	"github.com/go-telegram/bot"
	h "go_telegram_bot/src/handler"
	"os"
	"os/signal"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(h.Handler),
	}

	b, err := bot.New("5844620699:AAGbEPIFWKxTDr0jR_A77Rba95jtZBSQlGM", opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.CalendarHandler)

	b.Start(ctx)
}
