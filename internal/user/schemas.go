package user

type UserSchema struct {
	Id        uint
	UserName  string
	FirstName string
	Password  string `json:"-"`
}
