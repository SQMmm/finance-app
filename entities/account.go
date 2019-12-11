package entities

type Account struct {
	ID           int64
	User         *User
	Title        string
	StartBalance float64
	UseInReports bool
}
