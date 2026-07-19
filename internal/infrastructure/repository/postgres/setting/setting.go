package setting

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_telegram_bot/internal/infrastructure/repository/errors_rep"
	"go_telegram_bot/internal/infrastructure/repository/postgres"
	"go_telegram_bot/internal/usecase/booking"
	"time"
)

var _ booking.SettingRepository = (*SettingRepository)(nil)

var (
	settingList       map[string]interface{}
	SessionTimeHour   float64
	TimeKeyboarWidth  int
	DaysInSchedule    int
	StartHourSchedule map[string]float64
	EndHourSchedule   map[string]float64
)

type SettingRepository struct {
	DB *postgres.DB
}

func NewSettingRepository(db *postgres.DB) *SettingRepository {
	return &SettingRepository{
		DB: db,
	}
}

// GetDaysInSchedule
func (s *SettingRepository) GetDaysInSchedule(ctx context.Context) (float64, error) {
	return s.getNumberVal(ctx, "days_in_schedule")
}

// GetSessionTimeMinute
func (s *SettingRepository) GetSessionTimeMinute(ctx context.Context) (time.Duration, error) {
	sessionTime, err := s.getNumberVal(ctx, "session_time_minute")
	if err != nil {
		return 0, err
	}
	return time.Minute * time.Duration(sessionTime), nil
}

func (s *SettingRepository) stringToDuration(ctx context.Context, timeStr string) (time.Duration, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, nil
}

// GetDayStartTime
func (s *SettingRepository) GetDayStartTime(ctx context.Context, day string) (time.Duration, error) {
	var startHourSchedule map[string]string
	sched, err := s.GetJSONVal(ctx, "start_hour_schedule")
	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(sched, &startHourSchedule); err != nil {
		return 0, err
	}
	startTimeStr, ok := startHourSchedule[day]
	if !ok {
		return 0, errors_rep.ErrNullValue
	}
	startTime, err := s.stringToDuration(ctx, startTimeStr)
	if err != nil {
		return 0, err
	}
	return startTime, nil
}

// GetDayEndTime
func (s *SettingRepository) GetDayEndTime(ctx context.Context, day string) (time.Duration, error) {
	var startHourSchedule map[string]string
	sched, err := s.GetJSONVal(ctx, "end_hour_schedule")
	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(sched, &startHourSchedule); err != nil {
		return 0, err
	}
	startTimeStr, ok := startHourSchedule[day]
	if !ok {
		return 0, errors_rep.ErrNullValue
	}
	startTime, err := s.stringToDuration(ctx, startTimeStr)
	if err != nil {
		return 0, err
	}
	return startTime, nil
}

// getNumberVal
func (s *SettingRepository) getNumberVal(ctx context.Context, settingCode string) (float64, error) {
	var numberVal sql.NullFloat64
	query := "SELECT s.number_value FROM go_bot.setting s where s.setting_code = $1"
	row := s.DB.Pool.QueryRow(ctx, query, settingCode)
	err := row.Scan(&numberVal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors_rep.ErrNoRowsAffected
		}
		return 0, err
	}
	if !numberVal.Valid {
		return 0, errors_rep.ErrNullValue
	}
	return numberVal.Float64, nil
}

// GetStringVal
func (s *SettingRepository) GetStringVal(ctx context.Context, pool *pgxpool.Conn, setting_code string) (string, error) {
	var stringVal string
	query := "SELECT s.string_value FROM go_bot.setting s where s.setting_code = $1"
	row := s.DB.Pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(stringVal)
	if err != nil {
		return stringVal, err
	}
	return stringVal, nil
}

// GetDateVal
func (s *SettingRepository) GetDateVal(ctx context.Context, pool *pgxpool.Conn, setting_code string) (time.Time, error) {
	var dateVal time.Time
	query := "SELECT s.date_value FROM go_bot.setting s where s.setting_code = $1"
	row := pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(dateVal)
	if err != nil {
		return dateVal, err
	}
	return dateVal, nil
}

// GetJSONVal
func (s *SettingRepository) GetJSONVal(ctx context.Context, setting_code string) ([]byte, error) {
	var jsonVal []byte
	query := "SELECT s.json_value FROM go_bot.setting s where s.setting_code = $1"
	row := s.DB.Pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(&jsonVal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors_rep.ErrNoRowsAffected
		}
		return nil, err
	}
	if jsonVal == nil {
		return nil, errors_rep.ErrNullValue
	}
	return jsonVal, nil
}

//
//func LoadSettings(ctx context.Context) bool {
//	var (
//		err     error
//		jsonVal []byte
//		ok      bool = true
//	)
//	conn, err := db.PGPool.Acquire(ctx)
//	defer conn.Release()
//	if err != nil {
//		slog.Logger.Error("DB acquire error:", err.Error())
//	}
//
//	if SessionTimeHour, err = GetNumberVal(ctx, conn, "session_time_hour"); err != nil {
//		slog.Logger.Error("Error load session_time_hour:", err.Error())
//		ok = false
//	}
//	if val, err := GetNumberVal(ctx, conn, "time_keyboar_width"); err != nil {
//		slog.Logger.Error("Error load time_keyboar_width:", err.Error())
//		ok = false
//	} else {
//		TimeKeyboarWidth = int(val)
//	}
//	if val, err := GetNumberVal(ctx, conn, "days_in_schedule"); err != nil {
//		slog.Logger.Error("Error load days_in_schedule:", err.Error())
//		ok = false
//	} else {
//		DaysInSchedule = int(val)
//	}
//	if jsonVal, err = GetJSONVal(ctx, conn, "start_hour_schedule"); err != nil {
//		slog.Logger.Error("Error load start_hour_schedule:", err.Error())
//		ok = false
//	} else {
//		if err := json.Unmarshal(jsonVal, &StartHourSchedule); err != nil {
//			slog.Logger.Error("Error unmarshal start_hour_schedule:", err.Error())
//			ok = false
//		}
//	}
//	if jsonVal, err = GetJSONVal(ctx, conn, "end_hour_schedule"); err != nil {
//		slog.Logger.Error("Error load end_hour_schedule:", err.Error())
//		ok = false
//	} else {
//		if err := json.Unmarshal(jsonVal, &EndHourSchedule); err != nil {
//			slog.Logger.Error("Error unmarshal start_hour_schedule:", err.Error())
//			ok = false
//		}
//	}
//	return ok
//}
//
//func InitSettings() {
//	if ok := LoadSettings(context.Background()); !ok {
//		slog.Logger.Error("Load setting error!")
//	} else {
//		slog.Logger.Info("Load setting done!")
//	}
//	//fmt.Print(strconv.FormatFloat(SessionTimeHour, 'g', -1, 64))
//	//if val, err := json.Marshal(StartHourScheduler); err == nil {
//	//	fmt.Print(string(val))
//	//}
//	//if val, err := json.Marshal(EndHourScheduler); err == nil {
//	//	fmt.Print(string(val))
//	//}
//
//}
