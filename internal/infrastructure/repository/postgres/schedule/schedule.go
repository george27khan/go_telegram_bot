package schedule

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/infrastructure/repository/postgres"
	_ "go_telegram_bot/internal/usecase/booking"
	"time"
)

//var _ calendar.SettingRepository = (*ScheduleRepository)(nil)

type ScheduleRepository struct {
	DB *postgres.DB
}

func NewScheduleRepository(db *postgres.DB) *ScheduleRepository {
	return &ScheduleRepository{
		DB: db,
	}
}

func (sr *ScheduleRepository) Insert(ctx context.Context, schedule entity.Schedule) error {
	query := "INSERT INTO go_bot.schedule(id_user, id_employee, visit_dt) VALUES (@id_user, @id_employee, @visit_dt)"
	args := pgx.NamedArgs{
		"id_user":     schedule.ClientID,
		"id_employee": schedule.EmployeeID,
		"visit_dt":    schedule.VisitDt,
	}
	_, err := sr.DB.Pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (sr *ScheduleRepository) GetFreeEmpVisitDt(ctx context.Context, visitDt time.Time) (empSlice []entity.EmployeeID, err error) {
	var (
		EmployeeID entity.EmployeeID
	)

	query := "select e.id from go_bot.employee e where not exists(select 1 from go_bot.schedule t where t.id_employee = e.id and t.visit_dt = $1)"
	rows, errQuery := sr.DB.Pool.Query(ctx, query, visitDt)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&EmployeeID); err != nil {
			return nil, err
		}
		empSlice = append(empSlice, EmployeeID)
	}
	return
}

func (sr *ScheduleRepository) GetByDt(ctx context.Context, visitDt time.Time) (shedSlice []entity.Schedule, err error) {
	var (
		schedule entity.Schedule
	)
	query := "select id_user, visit_dt, created_dt, id_employee from go_bot.schedule t where t.id_employee = @id_employee and DATE_TRUNC('DAY', t.visit_dt) = $1"
	rows, err := sr.DB.Pool.Query(ctx, query, visitDt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&schedule.ClientID, &schedule.EmployeeID, &schedule.VisitDt); err != nil {
			return nil, err
		}
		shedSlice = append(shedSlice, schedule)
	}
	return
}

func (sr *ScheduleRepository) GetByClientID(ctx context.Context, idUser entity.ClientID) (entity.Schedule, error) {
	var (
		sched entity.Schedule
	)
	query := "select id_user, visit_dt, id_employee from go_bot.schedule t where t.id_user = $1 order by t.visit_dt desc limit 1"
	row := sr.DB.Pool.QueryRow(ctx, query, idUser)
	if err := row.Scan(&sched.ClientID, &sched.VisitDt, &sched.EmployeeID); err != nil {
		return entity.Schedule{}, err
	}
	return sched, nil
}

func (sr *ScheduleRepository) GetAllByUser(ctx context.Context, idUser int64) ([]entity.Schedule, error) {
	var (
		schedSlice []entity.Schedule
		sched      entity.Schedule
	)
	query := "select id_user, visit_dt, id_employee from go_bot.schedule t where t.id_user = $1 order by t.visit_dt desc"
	rows, err := sr.DB.Pool.Query(ctx, query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&sched.ClientID, &sched.VisitDt, &sched.EmployeeID); err != nil {
			return nil, err
		}
		schedSlice = append(schedSlice, sched)
	}
	return schedSlice, nil
}

func (sr *ScheduleRepository) IsFreeTime(ctx context.Context, visitDt time.Time) (ok bool, err error) {
	var cnt int
	query := "select count(1) from go_bot.schedule t where t.visit_dt = $1"
	row := sr.DB.Pool.QueryRow(ctx, query, visitDt)
	if err = row.Scan(&cnt); err != nil {
		return
	}
	if cnt > 0 {
		return false, nil
	}
	return true, nil
}

func (sr *ScheduleRepository) GetFreeEmployee(ctx context.Context, visitTime time.Time) ([]entity.Employee, error) {
	var (
		emp        entity.Employee
		idPosition int
		err        error
	)
	query := `select e.id, 
       				 e.first_name, 
       				 e.middle_name, 
       				 e.last_name, 
       				 e.birth_date, 
       				 e.email, 
       				 e.phone_number, 
       				 e.id_position, 
       				 e.hire_date, 
       				 e.photo,
       				 e.is_active
                from go_bot.employee e
               where not exists(select 1 
                                  from go_bot.schedule s 
                                 where s.visit_dt = $1
                                   and s.id_employee = e.id)`

	rows, err := sr.DB.Pool.Query(ctx, query, visitTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	empList := make([]entity.Employee, 0)
	for rows.Next() {
		err = rows.Scan(&emp.Id, &emp.FirstName, &emp.MiddleName, &emp.LastName, &emp.BirthDate, &emp.Email, &emp.PhoneNumber, &idPosition, &emp.HireDate, &emp.Photo, &emp.IsActive)
		if err != nil {
			return nil, err
		}
		empList = append(empList, emp)
	}

	return empList, nil
}

func (sr *ScheduleRepository) Save(ctx context.Context, schedule entity.Schedule) error {
	query := `insert into go_bot.schedule (id_user, id_employee, visit_dt)
              values ($1, $2, $3)`
	_, err := sr.DB.Pool.Exec(ctx, query, schedule.ClientID, schedule.EmployeeID, schedule.VisitDt)
	if err != nil {
		return err
	}
	return nil
}
