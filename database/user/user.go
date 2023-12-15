package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	db "go_telegram_bot/database"
)

func InsertUser(ctx context.Context, id int64, userName string, firstName string, lastName string, phone string) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return fmt.Errorf("Connection problem: %w", err)
	}
	query := "INSERT INTO go_bot.user(id, user_name, first_name, last_name, phone) VALUES (@id, @user_name, @first_name, @last_name, @phone)"

	args := pgx.NamedArgs{
		"id":         id,
		"user_name":  userName,
		"first_name": firstName,
		"last_name":  lastName,
		"phone":      phone,
	}
	_, err = pool.Exec(ctx, query, args)
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("Insert row problem: %w", err)
	}
	return nil
}
