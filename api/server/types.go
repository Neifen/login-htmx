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

func (u *UserType) ToUserReq() *userReq {
	token, exp, err := NewToken(u.uid, u.name)
	refresh, refreshExp := NewRefreshToken(u.uid)
	if err != nil {
		return emptyUser()
	}

	return &userReq{
		name:       u.name,
		uuid:       u.uid,
		isLoggedIn: true,
		token:      token, expires: exp,
		refresh:        refresh,
		refreshExpires: refreshExp,
	}
}
