package entities

type Category struct {
	ID    int64
	User  *User
	Title string
	Icon  *Icon
	Color *IconColor
}
