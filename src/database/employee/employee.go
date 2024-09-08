package employee

import (
	"context"
	"github.com/jackc/pgx/v5"
	db "go_telegram_bot/src/database"
	pstn "go_telegram_bot/src/database/position"
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
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
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
	if _, err := conn.Exec(ctx, query, args); err != nil {
		return err
	}

	return nil
}

// Delete функция для удаления записи из таблицы
func (e *Employee) Delete(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from go_bot.employee t where t.id = $1"
	_, err = conn.Exec(ctx, query, e.Id)
	if err := conn.Ping(ctx); err != nil {
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
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query := "select id, first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo from go_bot.employee t"
	rows, err := conn.Query(ctx, query)
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
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from go_bot.employee t where t.id = $1"
	if _, err = conn.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context, id int) (Employee, error) {
	var (
		emp        Employee
		idPosition int
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return Employee{}, err
	}
	query := "select id, first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo from go_bot.employee t where id=$1"
	row := conn.QueryRow(ctx, query, id)
	if err := row.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName, &emp.LastName, &emp.BirthDate, &emp.Email, &emp.PhoneNumber, &idPosition, &emp.HireDate, &emp.Photo); err != nil {
		return Employee{}, nil
	}
	if emp.Position, err = pstn.Get(ctx, idPosition); err != nil {
		return Employee{}, err
	}
	return emp, nil
}

func GetFIO(ctx context.Context, id int) (string, error) {
	var (
		FIO string
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return "", err
	}
	query := "select middle_name || ' ' || first_name || ' ' || last_name from go_bot.employee t where id=$1"
	row := conn.QueryRow(ctx, query, id)
	if err := row.Scan(&FIO); err != nil {
		return "", err
	}
	return FIO, nil
}
