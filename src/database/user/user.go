package user

import (
	"context"
	"github.com/jackc/pgx/v5"
	db "go_telegram_bot/src/database"
)

type User struct {
	Id        int64
	Name      string
	FirstName string
	LastName  string
	Phone     string
}

func Insert(ctx context.Context, id int64, userName string, firstName string, lastName string, phone string) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO go_bot.user(id, user_name, first_name, last_name, phone) VALUES (@id, @user_name, @first_name, @last_name, @phone)"

	args := pgx.NamedArgs{
		"id":         id,
		"user_name":  userName,
		"first_name": firstName,
		"last_name":  lastName,
		"phone":      phone,
	}
	_, err = conn.Exec(ctx, query, args)
	return nil
}

func IsExists(ctx context.Context, user_name string) (bool, error) {
	var userCnt int
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return false, err
	}
	query := "select count(id) from go_bot.user where user_name=$1"
	row := conn.QueryRow(ctx, query, user_name)
	if err := row.Scan(&userCnt); err != nil {
		return false, err
	}
	if userCnt > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func Get(ctx context.Context, id int) (User, error) {
	var (
		usr User
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return User{}, err
	}
	query := "select id, user_name, first_name, last_name, phone from go_bot.user where id=$1"
	row := conn.QueryRow(ctx, query, id)
	if err := row.Scan(&usr.Id, &usr.Name, &usr.FirstName, &usr.LastName, &usr.Phone); err != nil {
		return User{}, nil
	}
	return usr, nil

}
