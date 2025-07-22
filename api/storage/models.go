package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neifen/htmx-login/api/crypto"
)

type UserModel struct {
	id    int
	Uid   string
	Name  string
	Email string
	Pw    []byte
}

type RefreshTokenModel struct {
	id         int
	UserUid    string
	Token      string
	Expiration time.Time
	Remember   bool
}

func NewUserModel(name, email, pw string) *UserModel {

	pwHash, err := crypto.HashPassword(pw)
	if err != nil {
		fmt.Println("could not hash password")
		return nil
	}

	return &UserModel{
		Name:  name,
		Email: email,
		Pw:    pwHash,
		Uid:   uuid.NewString(),
	}
}

func NewRefreshTokenModel(userUid, token string, expiration time.Time, remember bool) *RefreshTokenModel {
	return &RefreshTokenModel{
		UserUid:    userUid,
		Token:      token,
		Expiration: expiration,
		Remember:   remember,
	}
}
