package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/api/logging"
	"github.com/neifen/htmx-login/view"
)

type HandlerSession struct {
	store Storage //interfaces are already pointers?
}

func NewHanderSession(store Storage) *HandlerSession {
	return &HandlerSession{
		store: store,
	}
}

func (s *HandlerSession) handleGetLogin(c echo.Context) error {
	if u := userFromToken(c); u.isLoggedIn {
		return s.redirectToHome(c, u)
	}

	child := view.Login()
	return view.RenderView(c, child)
}

func (s *HandlerSession) handlePostLogin(c echo.Context) error {
	if u := userFromToken(c); u.isLoggedIn {
		return s.redirectToHome(c, u)
	}

	email := c.FormValue("email")
	pw := c.FormValue("password")

	userReq := s.Authenticate(email, pw)

	/*
		Set-Cookie: access_token=eyJ…; HttpOnly; Secure
		Set-Cookie: refresh_token=…; Max-Age=31536000; Path=/api/auth/refresh; HttpOnly; Secure

	*/
	if userReq.isLoggedIn {
		//add cookie
		token := new(http.Cookie)
		token.Name = "token"
		token.Value = userReq.token
		token.Expires = *userReq.expires
		token.HttpOnly = true
		token.Secure = true
		c.SetCookie(token)
		tokenExp := new(http.Cookie)
		tokenExp.Name = "token-expires"
		tokenExp.Value = userReq.expires.String()
		c.SetCookie(tokenExp)

		refresh := new(http.Cookie)
		refresh.Name = "refresh"
		refresh.Path = "token/refresh"
		refresh.Value = userReq.refresh
		refresh.Expires = *userReq.refreshExpires
		refresh.HttpOnly = true
		refresh.Secure = true
		c.SetCookie(refresh)
		refreshExp := new(http.Cookie)
		refreshExp.Name = "refresh-expires"
		refreshExp.Value = userReq.refreshExpires.String()
		c.SetCookie(refreshExp)

		return s.redirectToHome(c, userReq)
	}

	// authenticate failed
	return s.redirectToLogin(c)
}

func (s *HandlerSession) handlePostLogout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = "delete"
	cookie.MaxAge = -1

	c.SetCookie(cookie)
	return s.redirectToLogin(c)
}

func (s *HandlerSession) handleGetRecovery(c echo.Context) error {
	if u := userFromToken(c); u.isLoggedIn {
		return s.redirectToHome(c, u)
	}

	child := view.PWRecovery()
	return view.RenderView(c, child)
}

func (s *HandlerSession) handleGetSignup(c echo.Context) error {
	if u := userFromToken(c); u.isLoggedIn {
		return s.redirectToHome(c, u)
	}

	child := view.Signup()
	return view.RenderView(c, child)
}

func (s *HandlerSession) handlePostSignup(c echo.Context) error {
	if u := userFromToken(c); u.isLoggedIn {
		return s.redirectToHome(c, u)
	}

	email := c.FormValue("email")
	pw := c.FormValue("password")
	name := c.FormValue("name")

	u := NewUserType(name, email, pw)
	err := s.store.CreateUser(u)

	if err != nil {
		logging.Error(err)
		// todo show error
		return s.handleGetSignup(c)
	}

	//todo success
	return s.redirectToLogin(c)
}

func (*HandlerSession) redirectToLogin(c echo.Context) error {
	child := view.Login()
	return view.ReplaceUrl(LOGIN_PATH, c, child)
}
