package server

import (
	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/view"
)

func (s *HandlerSession) handleGetHome(c echo.Context) error {
	u:= c.Get("u").(*userReq)
	//u := userFromToken(c)
	
	if !u.isLoggedIn {
		return s.redirectToLogin(c)
	}

	child := view.Home(u.name)
	return view.RenderView(c, child)
}

func (s *HandlerSession) redirectToHome(c echo.Context, user *userReq) error {
	child := view.Home(user.name)
	return view.ReplaceUrl(HOME_PATH, c, child)
}

func (s *HandlerSession) handleRandom(c echo.Context) error {
	// u:= c.Get("u").(*userReq)

	child := view.Random()
	return view.RenderView(c, child)
}
