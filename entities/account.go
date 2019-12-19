package entities

type Account struct {
	ID           int64
	User         *User
	Title        string
	StartBalance float64
	UseInReports bool
}

type Accounts []Account

type AccountGroup struct {
	ID       int64
	User     *User
	Title    string
	Accounts Accounts
}

func (as Accounts) GetIDs() []int64 {
	ids := make([]int64, len(as))
	for i, a := range as {
		ids[i] = a.ID
	}
	return ids
}
