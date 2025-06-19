package server

import (
	"fmt"
	"net/http"
	"time"

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
	if userReq.isLoggedIn {
		err := s.createAndHandleTokens(userReq, c)

		if err == nil {
			return s.redirectToHome(c, userReq)
		} else {
			fmt.Printf("\n\n auth: %v \n\n", err)

		}

	}

	// authenticate failed
	return s.redirectToLogin(c)
}

func (s *HandlerSession) handlePostTokenRefresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil || cookie == nil {
		return c.String(http.StatusUnauthorized, "No refresh token")
	}

	token, err := CheckRefreshToken(cookie.Value)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Refresh token invalid")
	}

	exp, err := token.GetExpiration()
	if err != nil || exp.Before(time.Now()) {
		return c.String(http.StatusUnauthorized, "Refresh token invalid")
	}

	//todo get token from db and check against this one
	refreshType, err := s.store.ReadRefreshTokenByToken(cookie.Value)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Refresh token invalid")
	}

	user, err := s.store.ReadUserByUid(refreshType.userUid)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Refresh token invalid")
	}

	err = s.createAndHandleTokens(user.ToUserReq(), c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Refresh token invalid")
	}

	return c.String(http.StatusOK, "Token successfully refreshed")
}

func (s *HandlerSession) createAndHandleTokens(user *userReq, c echo.Context) error {
	token, tokenExp, err := NewToken(user.uuid, user.name)
	if err != nil {
		return err
	}

	refresh, refreshExp, err := NewRefreshToken(user.uuid)
	if err != nil {
		return err
	}

	refreshType := NewRefreshTokenType(user.uuid, refresh, *refreshExp)
	err = s.store.CreateRefreshToken(refreshType)
	if err != nil {
		return err
	}

	/*
		Set-Cookie: access_token=eyJ…; HttpOnly; Secure
		Set-Cookie: refresh_token=…; Max-Age=31536000; Path=/api/auth/refresh; HttpOnly; Secure

	*/
	//add cookie
	tokenC := new(http.Cookie)
	tokenC.Name = "token"
	tokenC.Value = token
	tokenC.Expires = *tokenExp
	tokenC.HttpOnly = true
	tokenC.Secure = true
	c.SetCookie(tokenC)

	tokenExpC := new(http.Cookie)
	tokenExpC.Name = "token-expires"
	tokenExpC.Value = tokenExp.String()
	c.SetCookie(tokenExpC)

	refreshC := new(http.Cookie)
	refreshC.Name = "refresh"
	refreshC.Path = "token/refresh"
	refreshC.Value = refresh
	refreshC.Expires = *refreshExp
	refreshC.HttpOnly = true
	refreshC.Secure = true
	c.SetCookie(refreshC)

	refreshExpC := new(http.Cookie)
	refreshExpC.Name = "refresh-expires"
	refreshExpC.Value = refreshExp.String()
	c.SetCookie(refreshExpC)

	return nil
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
