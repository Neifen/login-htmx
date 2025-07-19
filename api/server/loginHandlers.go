package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/api/logging"
	"github.com/neifen/htmx-login/api/storage"
	"github.com/neifen/htmx-login/view"
)

type HandlerSession struct {
	store storage.Storage //interfaces are already pointers?
}

func NewHanderSession(store storage.Storage) *HandlerSession {
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

func (s *HandlerSession) handleTokenRefresh(c echo.Context) error {
	err := s.subHandleTokenRefresh(c)
	if err != nil {
		s.redirectToLogin(c)
	}

	return err
}

func (s *HandlerSession) subHandleTokenRefresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil || cookie == nil {
		return c.String(http.StatusUnauthorized, "no refresh token")
	}

	token, err := CheckRefreshToken(cookie.Value)
	if err != nil {
		return c.String(http.StatusUnauthorized, "refresh token invalid")
	}

	exp, err := token.GetExpiration()
	if err != nil || exp.Before(time.Now()) {
		return c.String(http.StatusUnauthorized, "refresh token expired")
	}

	//todo get token from db and check against this one
	refreshType, err := s.store.ReadRefreshTokenByToken(cookie.Value)
	if err != nil {
		return c.String(http.StatusUnauthorized, "could not load refresh token from db:"+err.Error())
	}

	fmt.Printf("Refresh Token for uid %s loaded\n", refreshType.UserUid)
	user, err := s.store.ReadUserByUid(refreshType.UserUid)
	if err != nil {
		return c.String(http.StatusUnauthorized, "user invalid")
	}

	err = s.createAndHandleTokens(userFromModel(user), c)
	if err != nil {
		return c.String(http.StatusUnauthorized, "creating new tokens failed")
	}

	//this is needed to set token before See Other
	// c.String(http.StatusOK, "Token successfully refreshed")

	returnUrl := c.QueryParam("return")
	if returnUrl != "" {
		return c.Redirect(http.StatusOK, returnUrl)
	}

	//todo: quicker redirect?
	return c.String(http.StatusOK, "Token successfully refreshed")
}

func redirectToTokenRefresh(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		// standard redirect
		return c.Redirect(http.StatusTemporaryRedirect, ("token/refresh?return=" + c.Request().URL.Path))
	}
	return c.String(http.StatusUnauthorized, "Unauthorized")
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

	refreshType := storage.NewRefreshTokenModel(user.uuid, refresh, *refreshExp)
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

	u := storage.NewUserModel(name, email, pw)
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
