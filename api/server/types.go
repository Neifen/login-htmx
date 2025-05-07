package server

type UserType struct {
	id    int
	name  string
	email string
	pw    string
}

func NewUserType(name, email, pw string) *UserType {
	return &UserType{
		name:  name,
		email: email,
		pw:    pw,
	}
}
