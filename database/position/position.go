package position

import (
	"context"
	"fmt"
	db "go_telegram_bot/database"
	"strings"
)

func InsertPosition(ctx context.Context, positionName string) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return fmt.Errorf("Connection problem: %w", err)
	}
	query := "INSERT INTO go_bot.position(position_name) VALUES ($1)"
	if _, err = pool.Exec(ctx, query, strings.ToLower(positionName)); err != nil {
		return fmt.Errorf("Insert row problem: %w", err)
	}
	return nil
}

func DeletePositionById(ctx context.Context, id int) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return fmt.Errorf("Connection problem: %w", err)
	}
	query := "delete from go_bot.position t where t.id = $1"
	if _, err = pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("Insert row problem: %w", err)
	}
	return nil
}

func SelectAllPosition(ctx context.Context) (positions []string, err error) {
	var position string
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return nil, fmt.Errorf("Connection problem: %w", err)
	}
	query := "select t.position_name from go_bot.position t"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&position); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		positions = append(positions, position)
	}
	return
}

func SelectAllPositionMap(ctx context.Context) (map[int]string, error) {
	var (
		position string
		id       int
	)
	positions := make(map[int]string)
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return nil, fmt.Errorf("Connection problem: %w", err)
	}
	query := "select t.id, t.position_name from go_bot.position t"

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &position); err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		positions[id] = position
	}
	return positions, nil
}
