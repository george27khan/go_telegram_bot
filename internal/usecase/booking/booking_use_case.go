package booking

import (
	"context"
	"go_telegram_bot/internal/domain/entity"
	"go_telegram_bot/internal/transport/telegram/handler"
	"time"
)

var _ handler.BookingUseCase = (*BookingUseCase)(nil)

type SettingRepository interface {
	GetDaysInSchedule(ctx context.Context) (float64, error)
	GetDayStartTime(ctx context.Context, day string) (time.Duration, error)
	GetDayEndTime(ctx context.Context, day string) (time.Duration, error)
	GetSessionTimeMinute(ctx context.Context) (time.Duration, error)
}

type SchedulerRepository interface {
	IsFreeTime(ctx context.Context, visitDt time.Time) (ok bool, err error)
	GetFreeEmployee(ctx context.Context, visitTime time.Time) ([]entity.Employee, error)
	Save(ctx context.Context, schedule entity.Schedule) error
}

type EmployeeRepository interface {
	Get(ctx context.Context, id int) (entity.Employee, error)
	GetAll(ctx context.Context) ([]entity.Employee, error)
}

type BookingUseCase struct {
	SettingRepository   SettingRepository
	SchedulerRepository SchedulerRepository
	EmployeeRepository  EmployeeRepository
}

func NewBookingUseCase(settingRepository SettingRepository, schedulerRepository SchedulerRepository, employeeRepository EmployeeRepository) *BookingUseCase {
	return &BookingUseCase{SettingRepository: settingRepository,
		SchedulerRepository: schedulerRepository,
		EmployeeRepository:  employeeRepository,
	}
}

func (b *BookingUseCase) GetCalendarDays(ctx context.Context) (dateFrom time.Time, dateTo time.Time, err error) {
	dateFrom = time.Now()
	daysInSchedule, err := b.SettingRepository.GetDaysInSchedule(ctx)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	dateTo = dateFrom.AddDate(0, 0, int(daysInSchedule))
	return
}

func (b *BookingUseCase) GetDayTimes(ctx context.Context, date time.Time) (timeFrom time.Time, timeTo time.Time, err error) {
	curDayName := time.Now().Format("Mon")
	startTime, err := b.SettingRepository.GetDayStartTime(ctx, curDayName)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime, err := b.SettingRepository.GetDayEndTime(ctx, curDayName)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	timeFrom = date.Add(startTime)
	timeTo = date.Add(endTime)
	return
}

func (b *BookingUseCase) GetSessionTimeMinute(ctx context.Context) (time.Duration, error) {
	sessionTime, err := b.SettingRepository.GetSessionTimeMinute(ctx)
	if err != nil {
		return 0, err
	}
	return sessionTime, nil
}

func (b *BookingUseCase) IsFreeTime(ctx context.Context, visitTime time.Time) (ok bool, err error) {
	ok, err = b.SchedulerRepository.IsFreeTime(ctx, visitTime)
	if err != nil {
		return
	}
	return
}

func (b *BookingUseCase) GetFreeEmployee(ctx context.Context, visitTime time.Time) ([]entity.Employee, error) {
	empListAll, err := b.SchedulerRepository.GetFreeEmployee(ctx, visitTime)
	if err != nil {
		return nil, err
	}
	return empListAll, nil
}

func (b *BookingUseCase) SaveSchedule(ctx context.Context, schedule entity.Schedule) error {
	if err := b.SchedulerRepository.Save(ctx, schedule); err != nil {
		return err
	}
	return nil
}
