package entity

const (
	StateMain     UserState = "Main"
	StateCalendar UserState = "Calendar"
	StateBanned   UserState = "banned"
)

type UserState string
