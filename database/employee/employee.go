package employee

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	db "go_telegram_bot/database"
	pstn "go_telegram_bot/database/position"
	"time"
)

// Employee тип для представления записи из таблицы employee
type Employee struct {
	Id          int
	FirstName   string
	MiddleName  string
	LastName    string
	BirthDate   time.Time
	Email       string
	PhoneNumber string
	Position    pstn.Position
	HireDate    time.Time
	Photo       []byte
}

// Insert функция для добавление записи в таблицу
func (e *Employee) Insert(ctx context.Context) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return err
	}
	query := "INSERT INTO go_bot.employee(first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo) VALUES (@first_name, @middle_name, @last_name, @birth_date, @email, @phone_number, @id_position, @hire_date, @photo)"

	args := pgx.NamedArgs{
		"first_name":   e.FirstName,
		"middle_name":  e.MiddleName,
		"last_name":    e.LastName,
		"birth_date":   e.BirthDate,
		"email":        e.Email,
		"phone_number": e.PhoneNumber,
		"id_position":  e.Position.Id,
		"hire_date":    e.HireDate,
		"photo":        e.Photo,
	}
	if res, err := pool.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println(res)
	}

	return nil
}

// Delete функция для удаления записи из таблицы
func (e *Employee) Delete(ctx context.Context) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return err
	}
	query := "delete from go_bot.employee t where t.id = $1"
	_, err = pool.Exec(ctx, query, e.Id)
	if err := pool.Ping(ctx); err != nil {
		return err
	}
	return nil
}

// SelectAll функция для получения всех записей из таблицы в виде среза элементов типа Employee
func SelectAll(ctx context.Context) ([]Employee, error) {
	var (
		employee  Employee
		employees []Employee
		idPos     int
	)
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return nil, err
	}
	query := "select id, first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo from go_bot.employee t"
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&employee.Id, &employee.FirstName, &employee.MiddleName, &employee.LastName, &employee.BirthDate, &employee.Email, &employee.PhoneNumber, &idPos, &employee.HireDate, &employee.Photo); err != nil {
			return nil, err
		}
		if employee.Position, err = pstn.Get(ctx, idPos); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func SelectAllStr(ctx context.Context) (res []string, err error) {
	emloyees, err := SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, empl := range emloyees {
		res = append(res, empl.MiddleName+" "+empl.FirstName+" "+empl.LastName)
	}
	return res, nil
}

// DeleteById удаление записи из таблицы по id
func DeleteById(ctx context.Context, id int) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return err
	}
	query := "delete from go_bot.employee t where t.id = $1"
	if _, err = pool.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}
