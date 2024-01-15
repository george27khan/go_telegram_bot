package setting

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	db "go_telegram_bot/database"
	"time"
)

var (
	settingList       map[string]interface{}
	SessionTimeHour   float64
	TimeKeyboarWidth  int
	DaysInSchedule    int
	StartHourSchedule map[string]float64
	EndHourSchedule   map[string]float64
)

func GetNumberVal(ctx context.Context, pool *pgxpool.Conn, setting_code string) (float64, error) {
	var numberVal float64
	query := "SELECT s.number_value FROM go_bot.setting s where s.setting_code = $1"
	row := pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(&numberVal)
	if err != nil {
		return numberVal, fmt.Errorf("unable to scan row: %w", err)
	}
	return numberVal, nil
}

func GetStringVal(ctx context.Context, pool *pgxpool.Conn, setting_code string) (string, error) {
	var stringVal string
	query := "SELECT s.string_value FROM go_bot.setting s where s.setting_code = $1"
	row := pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(stringVal)
	if err != nil {
		return stringVal, fmt.Errorf("unable to scan row: %w", err)
	}
	return stringVal, nil
}

func GetDateVal(ctx context.Context, pool *pgxpool.Conn, setting_code string) (time.Time, error) {
	var dateVal time.Time
	query := "SELECT s.date_value FROM go_bot.setting s where s.setting_code = $1"
	row := pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(dateVal)
	if err != nil {
		return dateVal, fmt.Errorf("unable to scan row: %w", err)
	}
	return dateVal, nil
}

func GetJSONVal(ctx context.Context, pool *pgxpool.Conn, setting_code string) ([]byte, error) {
	var (
		jsonVal []byte
	)
	query := "SELECT s.json_value FROM go_bot.setting s where s.setting_code = $1"
	row := pool.QueryRow(ctx, query, setting_code)
	err := row.Scan(&jsonVal)
	if err != nil {
		return nil, fmt.Errorf("unable to scan row: %w", err)
	}
	return jsonVal, nil
}

func LoadSettings(ctx context.Context) bool {
	var (
		err     error
		jsonVal []byte
		ok      bool = true
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()

	if SessionTimeHour, err = GetNumberVal(ctx, conn, "session_time_hour"); err != nil {
		fmt.Println("Error load session_time_hour")
		ok = false
	}
	if val, err := GetNumberVal(ctx, conn, "time_keyboar_width"); err != nil {
		fmt.Println("Error load time_keyboar_width")
		ok = false
	} else {
		TimeKeyboarWidth = int(val)
	}
	if val, err := GetNumberVal(ctx, conn, "days_in_schedule"); err != nil {
		fmt.Println("Error load days_in_schedule")
		ok = false
	} else {
		DaysInSchedule = int(val)
	}

	if jsonVal, err = GetJSONVal(ctx, conn, "start_hour_schedule"); err != nil {
		fmt.Println("Error load start_hour_schedule")
		ok = false
	} else {
		if err := json.Unmarshal(jsonVal, &StartHourSchedule); err != nil {
			fmt.Println("Error load start_hour_schedule")
		}
	}
	if jsonVal, err = GetJSONVal(ctx, conn, "end_hour_schedule"); err != nil {
		fmt.Println("Error load end_hour_schedule")
		ok = false
	} else {
		if err := json.Unmarshal(jsonVal, &EndHourSchedule); err != nil {
			fmt.Println("Error load start_hour_schedule")
		}
	}
	return ok
}

func init() {
	if ok := LoadSettings(context.Background()); !ok {
		fmt.Println("Load setting error!")
	}
	//fmt.Print(strconv.FormatFloat(SessionTimeHour, 'g', -1, 64))
	//if val, err := json.Marshal(StartHourScheduler); err == nil {
	//	fmt.Print(string(val))
	//}
	//if val, err := json.Marshal(EndHourScheduler); err == nil {
	//	fmt.Print(string(val))
	//}

}
