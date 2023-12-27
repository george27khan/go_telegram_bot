package main

import (
	"context"
	"github.com/go-telegram/bot"
	_ "go_telegram_bot/database/setting"
	h "go_telegram_bot/src/handler"
	"os"
	"os/signal"
)

const BotToken = "5844620699:AAGbEPIFWKxTDr0jR_A77Rba95jtZBSQlGM"

// Send any text message to the bot after the bot has been started
func main() {
	//db.DropDB()
	//db.InitDB()
	//res := rdb.HSet(context.Background(), "state", map[string]interface{}{"key1": "value1"})
	//fmt.Println("res ", res)
	//res1 := rdb.HGet(context.TODO(), "state", "key1")
	//fmt.Println("res1 ", res1.Val())

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
	//b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "/schedule", bot.MatchTypeExact, h.CalendarHandler)
	b.Start(ctx)

	//db_sett.LoadSettings(ctx)
}
