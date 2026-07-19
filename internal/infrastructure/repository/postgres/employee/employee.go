package employee

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go_telegram_bot/internal/domain/entity"
	errRep "go_telegram_bot/internal/infrastructure/repository/errors_rep"
	"go_telegram_bot/internal/infrastructure/repository/postgres"
)

type EmployeeRepository struct {
	DB *postgres.DB
}

func NewEmployeeRepository(db *postgres.DB) *EmployeeRepository {
	return &EmployeeRepository{
		DB: db,
	}
}

// Insert функция для добавление записи в таблицу
func (er *EmployeeRepository) Insert(ctx context.Context, employee entity.Employee) error {
	query := `INSERT INTO go_bot.employee(first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo) 
		     VALUES (@first_name, @middle_name, @last_name, @birth_date, @email, @phone_number, @id_position, @hire_date, @photo)`

	args := pgx.NamedArgs{
		"first_name":   employee.FirstName,
		"middle_name":  employee.MiddleName,
		"last_name":    employee.LastName,
		"birth_date":   employee.BirthDate,
		"email":        employee.Email,
		"phone_number": employee.PhoneNumber,
		"id_position":  employee.Position,
		"hire_date":    employee.HireDate,
		"photo":        employee.Photo,
	}
	if _, err := er.DB.Pool.Exec(ctx, query, args); err != nil {
		return err
	}
	return nil
}

// Delete функция для удаления записи из таблицы
func (er *EmployeeRepository) Delete(ctx context.Context, id entity.EmployeeID) error {
	query := "delete from go_bot.employee t where t.id = $1"
	cmdTag, err := er.DB.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errRep.ErrNoRowsAffected
	}
	return nil
}

// SelectAll функция для получения всех записей из таблицы в виде среза элементов типа Employee
func (er *EmployeeRepository) SelectAll(ctx context.Context) ([]entity.Employee, error) {
	var (
		employee  entity.Employee
		employees []entity.Employee
		idPos     int
	)

	query := "select id, first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo from go_bot.employee t"
	rows, err := er.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&employee.Id, &employee.FirstName, &employee.MiddleName, &employee.LastName, &employee.BirthDate, &employee.Email, &employee.PhoneNumber, &idPos, &employee.HireDate, &employee.Photo); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (er *EmployeeRepository) SelectAllStr(ctx context.Context) (res []string, err error) {
	emloyees, err := er.SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, empl := range emloyees {
		res = append(res, empl.MiddleName+" "+empl.FirstName+" "+empl.LastName)
	}
	return res, nil
}

// DeleteById удаление записи из таблицы по id
func (er *EmployeeRepository) DeleteById(ctx context.Context, id int) error {
	query := "delete from go_bot.employee t where t.id = $1"
	if _, err := er.DB.Pool.Exec(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (er *EmployeeRepository) Get(ctx context.Context, id int) (entity.Employee, error) {
	var (
		emp        entity.Employee
		idPosition int
		err        error
	)
	query := "select id, first_name, middle_name, last_name, birth_date, email, phone_number, id_position, hire_date, photo from go_bot.employee t where id=$1"
	row := er.DB.Pool.QueryRow(ctx, query, id)
	if err := row.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName, &emp.LastName, &emp.BirthDate, &emp.Email, &emp.PhoneNumber, &idPosition, &emp.HireDate, &emp.Photo); err != nil {
		return entity.Employee{}, nil
	}
	if err != nil {
		return entity.Employee{}, err
	}
	return emp, nil
}

func (er *EmployeeRepository) GetAll(ctx context.Context) ([]entity.Employee, error) {
	var (
		emp        entity.Employee
		idPosition int
		err        error
	)
	query := `select id, 
       				 first_name, 
       				 middle_name, 
       				 last_name, 
       				 birth_date, 
       				 email, 
       				 phone_number, 
       				 id_position, 
       				 hire_date, 
       				 photo,
       				 is_active
                from go_bot.employee`
	rows, err := er.DB.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	empList := make([]entity.Employee, 0)
	for rows.Next() {
		err = rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName, &emp.LastName, &emp.BirthDate, &emp.Email, &emp.PhoneNumber, &idPosition, &emp.HireDate, &emp.Photo)
		if err != nil {
			return nil, err
		}
		empList = append(empList, emp)
	}

	return empList, nil
}

func (er *EmployeeRepository) GetFIO(ctx context.Context, id int) (string, error) {
	var (
		FIO string
	)
	query := "select middle_name || ' ' || first_name || ' ' || last_name from go_bot.employee t where id=$1"
	row := er.DB.Pool.QueryRow(ctx, query, id)
	if err := row.Scan(&FIO); err != nil {
		return "", err
	}
	return FIO, nil
}
