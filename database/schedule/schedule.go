package schedule

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	db "go_telegram_bot/database"
	"time"
)

func InsertSchedule(ctx context.Context, idUser int64, visitDt time.Time) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return fmt.Errorf("Connection problem: %w", err)
	}
	query := "INSERT INTO go_bot.schedule(id_user, visit_dt) VALUES (@id_user, @visit_dt)"
	//query := "INSERT INTO go_bot.schedule(id_user, visit_dt) VALUES (1, )"
	args := pgx.NamedArgs{
		"id_user":  idUser,
		"visit_dt": visitDt,
	}
	res, err := pool.Exec(ctx, query, args)
	fmt.Println(res, idUser, visitDt.String(), err)
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("Insert row problem: %w", err)
	}
	return nil
}
