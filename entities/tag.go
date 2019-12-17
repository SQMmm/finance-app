package entities

type Tag struct {
	ID    int64
	User  *User
	Title string
}
