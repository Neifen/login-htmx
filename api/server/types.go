package server

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserType struct {
	id    int
	uid   string
	name  string
	email string
	pw    []byte
}

type RefreshTokenType struct {
	id         int
	userUid    string
	token      string
	expiration time.Time
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

func NewRefreshTokenType(userUid, token string, expiration time.Time) *RefreshTokenType {

	return &RefreshTokenType{
		userUid:    userUid,
		token:      token,
		expiration: expiration,
	}
}

func (u *UserType) ToUserReq() *userReq {
	return &userReq{
		name:       u.name,
		uuid:       u.uid,
		isLoggedIn: true,
	}
}
