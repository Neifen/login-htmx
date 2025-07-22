package server

import (
	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/api/crypto"
	"github.com/neifen/htmx-login/api/storage"
)

type userReq struct {
	isLoggedIn bool
	name       string
	uuid       string
}

func emptyUser() *userReq {
	return new(userReq)
}

func userFromModel(u *storage.UserModel) *userReq {
	return &userReq{
		name:       u.Name,
		uuid:       u.Uid,
		isLoggedIn: true,
	}
}

func userFromToken(c echo.Context) *userReq {
	cookie, err := c.Cookie("token")
	if err != nil || cookie == nil {
		return emptyUser()
	}

	token, err := crypto.ValidTokenFromCookies(cookie)
	if err != nil {
		return emptyUser()
	}

	uid, err := token.UserID()
	if err != nil {
		return emptyUser()
	}

	name, _ := token.UserName()

	return &userReq{isLoggedIn: true, name: name, uuid: uid}
}
