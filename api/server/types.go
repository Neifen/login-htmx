package server

import (
	"fmt"

	"github.com/google/uuid"
)

type UserType struct {
	id    int
	uid   string
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
		uid:   uuid.NewString(),
	}
}

func (u *UserType) ToUserInfo() *UserInfo {
	return &UserInfo{userName: u.name, uid: u.uid, crypt: nil}
}
