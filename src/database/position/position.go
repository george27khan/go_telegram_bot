package position

import (
	"context"
	db "go_telegram_bot/src/database"
	"strings"
)

// Position тип для представления записи из таблицы position
type Position struct {
	Id           int
	PositionName string
}

// Get функция для получения записи из position по id
func Get(ctx context.Context, id int) (Position, error) {
	var (
		position Position
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return Position{}, err
	}
	query := "select t.id, t.position_name from go_bot.position t where t.id=$1"

	row := conn.QueryRow(ctx, query, id)
	if err != nil {
		return Position{}, err
	}
	if err := row.Scan(&position.Id, &position.PositionName); err != nil {
		return Position{}, err
	}
	return position, nil
}

// GetByName функция для получения записи из position по наименованию
func GetByName(ctx context.Context, name string) (Position, error) {
	var (
		position Position
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return Position{}, err
	}
	query := "select t.id, t.position_name from go_bot.position t where t.position_name=$1"

	row := conn.QueryRow(ctx, query, name)
	if err != nil {
		return Position{}, err
	}
	if err := row.Scan(&position.Id, &position.PositionName); err != nil {
		return Position{}, err
	}
	return position, nil
}

// Insert функция для добавления записи в таблицу position
func (p *Position) Insert(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO go_bot.position(position_name) VALUES ($1)"
	if _, err := conn.Exec(ctx, query, strings.ToLower(p.PositionName)); err != nil {
		return err
	}
	return nil
}

// DeleteById функция для удаления записи из таблицы position по id
func DeleteById(ctx context.Context, id int) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from go_bot.position t where t.id = $1"
	if _, err = conn.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

// SelectAll функция для получения всех записей из таблицы position в виде среза элементов типа Position
func SelectAll(ctx context.Context) ([]Position, error) {
	var (
		position  Position
		positions []Position
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query := "select t.id, t.position_name from go_bot.position t"

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&position.Id, &position.PositionName); err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	return positions, nil
}

// SelectAllStr функция для получения всех записей из таблицы position в виде среза элементов типа string
func SelectAllStr(ctx context.Context) (res []string, err error) {
	positions, err := SelectAll(ctx)
	for _, position := range positions {
		res = append(res, position.PositionName)
	}
	return
}
