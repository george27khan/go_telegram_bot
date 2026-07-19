package client

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/infrastructure/repository/postgres"
	"go_telegram_bot/internal/usecase/start"
)

var _ start.ClientRepository = (*ClientRepository)(nil)

type ClientRepository struct {
	DB *postgres.DB
}

func NewClientRepository(db *postgres.DB) *ClientRepository {
	return &ClientRepository{
		DB: db,
	}
}

func (cr *ClientRepository) Create(ctx context.Context, client entity.Client) error {
	query := `INSERT INTO go_bot.user(id, user_name, first_name, last_name, phone) 
              VALUES (@id, @user_name, @first_name, @last_name, @phone)`

	args := pgx.NamedArgs{
		"id":         client.ID,
		"user_name":  client.UserName,
		"first_name": client.FirstName,
		"last_name":  client.LastName,
		"phone":      client.Phone,
	}
	_, err := cr.DB.Pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ClientRepository) IsExists(ctx context.Context, clientID entity.ClientID) (bool, error) {
	var userCnt int
	query := "select count(1) from go_bot.user where id=$1"
	row := cr.DB.Pool.QueryRow(ctx, query, clientID)
	if err := row.Scan(&userCnt); err != nil {
		return false, err
	}
	if userCnt > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (cr *ClientRepository) Get(ctx context.Context, id int) (entity.Client, error) {
	var (
		usr entity.Client
	)
	query := "select id, user_name, first_name, last_name, phone from go_bot.user where id=$1"
	row := cr.DB.Pool.QueryRow(ctx, query, id)
	if err := row.Scan(&usr.ID, &usr.UserName, &usr.FirstName, &usr.LastName, &usr.Phone); err != nil {
		return entity.Client{}, nil
	}
	return usr, nil

}
