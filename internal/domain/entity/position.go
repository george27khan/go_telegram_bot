package entity

type PositionID int

// Position тип для представления записи из таблицы position
type Position struct {
	Id           PositionID
	PositionName string
}
