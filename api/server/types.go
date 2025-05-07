package server

import "fmt"

type UserType struct {
	id    int
	name  string
	email string
	pw    []byte
}

func NewUserType(name, email, pw string) *UserType {

	pwHash, err := HashPassword(pw)
	if err != nil {
		fmt.Println("could not hash password")
		return nil
	}

	return &UserType{
		name:  name,
		email: email,
		pw:    pwHash,
	}
}
