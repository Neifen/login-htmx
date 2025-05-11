package server

import (
	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/view"
)

func (s *HandlerSession) handleGetHome(c echo.Context) error {
	// fmt.Println(c.Cookies())
	if !s.isLoggedIn(c.Cookie("token")) {
		return s.redirectToLogin(c)
	}

	child := view.Home(s.user.userName)
	return view.RenderView(c, child)
}

func (s *HandlerSession) redirectToHome(c echo.Context) error {
	child := view.Home(s.user.userName)
	return view.ReplaceUrl(HOME_PATH, c, child)
}
