package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	db "go_telegram_bot/src/database"
	stng "go_telegram_bot/src/database/setting"
	h "go_telegram_bot/src/handler"
	"log"
	"os"
	"os/signal"
)

var (
	botToken string
)

func loadEnv() {
	// loads DB settings from .env into the system
	if err := godotenv.Load("./.env"); err != nil {
		log.Print("No .env file found")
	}
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
}

// Send any text message to the bot after the bot has been started
func main() {
	loadEnv()
	db.InitDB() // разворачивание миграции при первом запуске, костыль
	db.ConnectDB()
	stng.InitSettings()
	//db.DropDB()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	//
	opts := []bot.Option{
		bot.WithDefaultHandler(h.DefaultHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartHandler)
	b.Start(ctx)
}
