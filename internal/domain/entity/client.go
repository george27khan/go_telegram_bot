package entity

type (
	ClientID int64
)

type Client struct {
	ID           ClientID
	UserName     string
	FirstName    string
	LastName     string
	LanguageCode string
	Phone        string
}
