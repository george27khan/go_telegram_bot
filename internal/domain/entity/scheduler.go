package entity

import "time"

type Schedule struct {
	ClientID   ClientID
	EmployeeID EmployeeID
	VisitDt    time.Time
}
