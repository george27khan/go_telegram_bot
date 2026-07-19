package position

import (
	"context"
	"go_telegram_bot/internal/domain/entity"
	errRep "go_telegram_bot/internal/infrastructure/repository/errors_rep"
	"go_telegram_bot/internal/infrastructure/repository/postgres"
	"strings"
)

type PositionRepository struct {
	DB *postgres.DB
}

func NewPositionRepository(db *postgres.DB) *PositionRepository {
	return &PositionRepository{
		DB: db,
	}
}

// Get функция для получения записи из position по id
func (pr *PositionRepository) Get(ctx context.Context, id int) (entity.Position, error) {
	var position entity.Position
	query := "select t.id, t.position_name from go_bot.position t where t.id=$1"

	row := pr.DB.Pool.QueryRow(ctx, query, id)
	if err := row.Scan(&position.Id, &position.PositionName); err != nil {
		return entity.Position{}, err
	}

	return position, nil
}

// GetByName функция для получения записи из position по наименованию
func (pr *PositionRepository) GetByName(ctx context.Context, name string) (entity.Position, error) {
	var position entity.Position

	query := "select t.id, t.position_name from go_bot.position t where t.position_name=$1"

	row := pr.DB.Pool.QueryRow(ctx, query, name)
	if err := row.Scan(&position.Id, &position.PositionName); err != nil {
		return entity.Position{}, err
	}
	return position, nil
}

// Insert функция для добавления записи в таблицу position
func (pr *PositionRepository) Insert(ctx context.Context, p entity.Position) error {
	query := "INSERT INTO go_bot.position(position_name) VALUES ($1)"
	cmdTag, err := pr.DB.Pool.Exec(ctx, query, strings.ToLower(p.PositionName))
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errRep.ErrNoRowsAffected
	}
	return nil
}

// DeleteById функция для удаления записи из таблицы position по id
func (pr *PositionRepository) DeleteById(ctx context.Context, id int) error {
	query := "delete from go_bot.position t where t.id = $1"
	if _, err := pr.DB.Pool.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

// SelectAll функция для получения всех записей из таблицы position в виде среза элементов типа Position
func (pr *PositionRepository) SelectAll(ctx context.Context) ([]entity.Position, error) {
	var (
		position  entity.Position
		positions []entity.Position
	)

	query := "select t.id, t.position_name from go_bot.position t"

	rows, err := pr.DB.Pool.Query(ctx, query)
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
func (pr *PositionRepository) SelectAllStr(ctx context.Context) (res []string, err error) {
	positions, err := pr.SelectAll(ctx)
	for _, position := range positions {
		res = append(res, position.PositionName)
	}
	return
}
