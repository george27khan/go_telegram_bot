package app

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/infrastructure/repository/postgres"
	repCli "go_telegram_bot/internal/infrastructure/repository/postgres/client"
	repEmp "go_telegram_bot/internal/infrastructure/repository/postgres/employee"
	repSched "go_telegram_bot/internal/infrastructure/repository/postgres/schedule"
	repSett "go_telegram_bot/internal/infrastructure/repository/postgres/setting"
	"go_telegram_bot/internal/pkg/logger"
	"go_telegram_bot/internal/transport/telegram/handler"
	router2 "go_telegram_bot/internal/transport/telegram/router"
	"go_telegram_bot/internal/usecase/booking"
	"go_telegram_bot/internal/usecase/setting"
	"go_telegram_bot/internal/usecase/start"
	"go_telegram_bot/internal/usecase/state"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func Run() {
	loadEnv()
	slog := logger.InitLogging()
	//глобальный контекст для отмены фоновых загрузок при остановке приложения
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	//
	pool, err := postgres.NewPostgresPool(rootCtx)
	if err != nil {
		log.Fatal(err)
	}
	r := postgres.NewRepository(pool, slog)
	clientRep := repCli.NewClientRepository(r)
	settingRep := repSett.NewSettingRepository(r)
	scheduleRep := repSched.NewScheduleRepository(r)
	empRep := repEmp.NewEmployeeRepository(r)
	startUC := start.NewStartUseCase(clientRep)
	bookingUC := booking.NewBookingUseCase(settingRep, scheduleRep, empRep)
	stateUC := state.NewStateUseCase()
	settingUC := setting.NewSettingUseCase(settingRep)

	bookingH := handler.NewBookingHandler(bookingUC, slog)
	settingH := handler.NewSettingHandler(settingUC, slog)

	router := router2.NewRouter()
	router.Register(entity.ActionCalendar, bookingH.GetCalendar)
	router.Register(entity.ActionSetting, settingH.GetSetting)

	startH := handler.NewStartHandler(startUC, stateUC, router, slog)
	opts := []bot.Option{
		//bot.WithDefaultHandler(h.DefaultHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startH.Start)
	b.Start(rootCtx)

	//gracefull shutdown
	<-rootCtx.Done() // ожидание сигнала завершения
	slog.Info("Start gracefull shutdown ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := b.Close(ctx); err != nil {
		slog.Error(fmt.Sprintf("Server Shutdown err: %s", err))
	}
	slog.Info("Server exiting.")
	pool.Close() // закрываем пул к базе
	slog.Info("Postgres pool closed.")
}
