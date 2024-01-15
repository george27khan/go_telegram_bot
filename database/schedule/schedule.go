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
	IdEmployee int
	VisitDt    time.Time
}

func (sched *Schedule) Insert(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO go_bot.schedule(id_user, id_employee, visit_dt) VALUES (@id_user, @id_employee, @visit_dt)"
	args := pgx.NamedArgs{
		"id_user":     sched.IdUser,
		"id_employee": sched.IdEmployee,
		"visit_dt":    sched.VisitDt,
	}
	_, err = conn.Exec(ctx, query, args)
	return nil
}

func GetFreeEmpVisitDt(ctx context.Context, visitDt time.Time) (empSlice []emp.Employee, err error) {
	var (
		idEmp int
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query := "select e.id from go_bot.employee e where not exists(select 1 from go_bot.schedule t where t.id_employee = e.id and t.visit_dt = $1)"
	rows, errQuery := conn.Query(ctx, query, visitDt)
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

func GetByDt(ctx context.Context, visitDt time.Time) (shedSlice []Schedule, err error) {
	var (
		schedule Schedule
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query := "select id_user, visit_dt, created_dt, id_employee from go_bot.schedule t where t.id_employee = @id_employee and DATE_TRUNC('DAY', t.visit_dt) = $1"
	rows, err := conn.Query(ctx, query, visitDt)
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

func GetByUser(ctx context.Context, idUser int64) (Schedule, error) {
	var (
		sched Schedule
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return Schedule{}, err
	}
	query := "select id_user, visit_dt, id_employee from go_bot.schedule t where t.id_user = $1 order by t.visit_dt desc limit 1"
	row := conn.QueryRow(ctx, query, idUser)
	if err != nil {
		return Schedule{}, err
	}
	if err := row.Scan(&sched.IdUser, &sched.VisitDt, &sched.IdEmployee); err != nil {
		return Schedule{}, err
	}
	return sched, nil
}

func GetAllByUser(ctx context.Context, idUser int64) ([]Schedule, error) {
	var (
		schedSlice []Schedule
		sched      Schedule
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	query := "select id_user, visit_dt, id_employee from go_bot.schedule t where t.id_user = $1 order by t.visit_dt desc"
	rows, err := conn.Query(ctx, query, idUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&sched.IdUser, &sched.VisitDt, &sched.IdEmployee); err != nil {
			return nil, err
		}
		schedSlice = append(schedSlice, sched)
	}
	return schedSlice, nil
}

func TimeExists(ctx context.Context, visitDt time.Time) (cnt int) {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return 0
	}
	query := "select count(1) from go_bot.schedule t where t.visit_dt = $1"
	row := conn.QueryRow(ctx, query, visitDt)
	if err != nil {
		return 1
	}
	if err := row.Scan(&cnt); err != nil {
		fmt.Println(err)
		return 1
	}
	return cnt

}
