package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"go_telegram_bot/src/slog"
	"os"
	"sync"
)

var PGPool *pgxpool.Pool

func getPGconnStr() (connStr string) {
	// loads DB settings from .env into the system
	//if err := godotenv.Load("./db.env"); err != nil {
	//	slog.Logger.Error("File ./db.env not found")
	//}
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")

	//"postgres://username:password@localhost:5432/database_name"
	connStr = fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPwd, host, dbName)
	slog.Logger.Info("DB Connection ", connStr)
	return
}

var (
	pgOnce sync.Once
)

func Pool(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, getPGconnStr())
}

func connect(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, getPGconnStr())
	if err != nil {
		slog.Logger.Error("Unable to connect to database:", err.Error())
		os.Exit(1)
	}
	return conn
}

func getMigrator(ctx context.Context, conn *pgx.Conn) *migrate.Migrator {
	migrator, err := migrate.NewMigrator(ctx, conn, "go_bot")
	if err != nil {
		slog.Logger.Error("Unable to create a migrator:", err.Error())
	}

	err = migrator.LoadMigrations(os.DirFS("./database/migration"))
	if err != nil {
		slog.Logger.Error("Unable to load migrations:", err.Error())
	}
	return migrator
}

func InitDB() {
	ctx := context.Background()
	conn := connect(ctx)
	defer conn.Close(ctx)

	migrator := getMigrator(ctx, conn)

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		slog.Logger.Error("Unable to get current schema version:", err.Error())
	}
	if ver == 0 {
		err := migrator.Migrate(ctx)
		if err != nil {
			slog.Logger.Error("Unable to migrate:", err.Error())
		}

		ver, err := migrator.GetCurrentVersion(ctx)
		if err != nil {
			slog.Logger.Error("Unable to get current schema version:", err.Error())
		}
		slog.Logger.Info("Migration done. Current schema version:", ver)
	}
}

func DropDB() {
	ctx := context.Background()
	conn := connect(ctx)
	defer conn.Close(ctx)

	migrator := getMigrator(ctx, conn)

	err := migrator.MigrateTo(ctx, 0)
	if err != nil {
		slog.Logger.Error("Unable to migrate:", err.Error())
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		slog.Logger.Error("Unable to get current schema version:", err.Error())
	}
	slog.Logger.Info("Migration done. Current schema version:", ver)
}

func ConnectDB() {
	var err error
	PGPool, err = Pool(context.Background())
	if err != nil {
		slog.Logger.Error("DB connection error:", err.Error())
	}
}
