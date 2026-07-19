package entity

import (
	//pstn "go_telegram_bot/internal/infrastructure/repository/postgres/position"
	"time"
)

type (
	EmployeeID int
)

// Employee структура для записи из таблицы employee
type Employee struct {
	Id          EmployeeID
	FirstName   string
	MiddleName  string
	LastName    string
	BirthDate   time.Time
	Email       string
	PhoneNumber string
	Position    Position
	HireDate    time.Time
	Photo       []byte
	IsActive    bool
}
