package schedule

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	db "go_telegram_bot/database"
	emp "go_telegram_bot/database/employee"
	"time"
)

type Schedule struct {
	IdUser     int64
	IdEmployee int64
	VisitDt    time.Time
}

func (sched *Schedule) Insert(ctx context.Context) error {
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return err
	}
	query := "INSERT INTO go_bot.schedule(id_user, id_employee, visit_dt) VALUES (@id_user, @id_employee, @visit_dt)"
	args := pgx.NamedArgs{
		"id_user":     sched.IdUser,
		"id_employee": sched.IdEmployee,
		"visit_dt":    sched.VisitDt,
	}
	_, err = pool.Exec(ctx, query, args)
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("Insert row problem: %w", err)
	}
	return nil
}

func GetFreeEmpVisitDt(ctx context.Context, visitDt time.Time) (empSlice []emp.Employee, err error) {
	var (
		idEmp int
	)
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return nil, err
	}
	query := "select e.id from go_bot.employee e where not exists(select 1 from go_bot.schedule t where t.id_employee = e.id and t.visit_dt = $1)"
	rows, errQuery := pool.Query(ctx, query, visitDt)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&idEmp); err != nil {
			return nil, err
		}
		if empl, errEmp := emp.Get(ctx, idEmp); err != nil {
			return nil, errEmp
		} else {
			empSlice = append(empSlice, empl)
		}
	}
	return empSlice, nil
}

func LoadSchedByEmpIdDt(ctx context.Context, visitDt time.Time) (shedSlice []Schedule, err error) {
	var (
		schedule Schedule
	)
	pool, err := db.Pool(ctx)
	defer pool.Close()
	if err != nil {
		return nil, err
	}
	query := "select id_user, visit_dt, created_dt, id_employee from go_bot.schedule t where t.id_employee = @id_employee and DATE_TRUNC('DAY', t.visit_dt) = $1"
	rows, err := pool.Query(ctx, query, visitDt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&schedule.IdUser, &schedule.IdEmployee, &schedule.VisitDt); err != nil {
			return nil, err
		}
		shedSlice = append(shedSlice, schedule)
	}
	return
}
