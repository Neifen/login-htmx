package server

import (
	"github.com/labstack/echo/v4"
)

type userReq struct {
	isLoggedIn     bool
	name           string
	uuid           string
}

func emptyUser() *userReq {
	return new(userReq)
}

func userFromToken(c echo.Context) *userReq {
	cookie, err := c.Cookie("token")
	if err != nil || cookie == nil {
		return emptyUser()
	}

	token, err := CheckToken(cookie.Value)
	if err != nil {
		return emptyUser()
	}

	uid, err := token.GetString("user-id")
	if err != nil {
		return emptyUser()
	}

	name, _ := token.GetString("user-name")

	return &userReq{isLoggedIn: true, name: name, uuid: uid}
}
