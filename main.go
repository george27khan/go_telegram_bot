package main

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
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
	if err := godotenv.Load("./bot.env"); err != nil {
		log.Print("No .env file found")
	}
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
}

const BotToken = "5844620699:AAGbEPIFWKxTDr0jR_A77Rba95jtZBSQlGM"

// Send any text message to the bot after the bot has been started
func main() {
	loadEnv()
	//db.DropDB()
	//db.InitDB()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	//
	opts := []bot.Option{
		bot.WithDefaultHandler(h.DefaultHandler),
	}

	b, err := bot.New(BotToken, opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.StartHandler)
	b.Start(ctx)
}
