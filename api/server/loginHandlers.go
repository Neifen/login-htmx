package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/api/crypto"
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
	remember := c.FormValue("remember") == "on"

	userReq := s.Authenticate(email, pw)
	if userReq.isLoggedIn {
		err := s.createAndHandleTokens(userReq, c, remember)

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
		fmt.Println(err)
		s.redirectToLogin(c)
	}

	return err
}

func (s *HandlerSession) subHandleTokenRefresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh")
	if err != nil {
		return errors.Wrapf(err, "getting refresh token from cookie %q failed", "refresh")
	}

	if cookie == nil {
		return fmt.Errorf("no refresh token in cookie %q", "refresh")
	}

	token, err := crypto.ValidTokenFromCookies(cookie)
	if err != nil {
		err = s.store.DeleteRefreshTokenByToken(token.Encrypted)
		if err != nil {
			fmt.Printf("could not delete refresh token from db %v\n", err)
		}

		return fmt.Errorf("refresh token could not be validated")
	}

	exp := token.Expiration
	fmt.Printf("cookie expires: %v, token expires: %v\n", cookie.Expires, token.Expiration)

	if exp.Before(time.Now()) {
		err = s.store.DeleteRefreshTokenByToken(token.Encrypted)
		if err != nil {
			fmt.Printf("could not delete refresh token from db %v\n", err)
		}
		return fmt.Errorf("refresh token expired")
	}

	refreshType, err := s.store.ReadRefreshTokenByToken(token.Encrypted)
	if err != nil {
		return errors.Wrapf(err, "could not load refresh token from db")
	}

	user, err := s.store.ReadUserByUid(refreshType.UserUid)
	if err != nil {
		return errors.Wrapf(err, "user invalid")
	}

	err = s.createAndHandleTokens(userFromModel(user), c, refreshType.Remember)
	if err != nil {
		return errors.Wrapf(err, "creating new tokens failed")
	}

	err = s.store.DeleteRefreshToken(refreshType)
	if err != nil {
		return errors.Wrapf(err, "could not delete old refresh token")
	}

	returnUrl := c.QueryParam("return")
	if returnUrl != "" {
		fmt.Printf("redirect with return: %s \n", returnUrl)
		return c.Redirect(http.StatusSeeOther, returnUrl)
	}

	fmt.Printf("redirect with no return\n")
	return c.String(http.StatusOK, "Token successfully refreshed")
}

func redirectToTokenRefresh(c echo.Context) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		// standard redirect
		return c.Redirect(http.StatusTemporaryRedirect, ("token/refresh?return=" + c.Request().URL.Path))
	}
	return c.String(http.StatusUnauthorized, "Unauthorized")
}

func (s *HandlerSession) createAndHandleTokens(user *userReq, c echo.Context, remember bool) error {
	access, err := crypto.NewAccessToken(user.uuid, user.name)
	if err != nil {
		return errors.Wrap(err, "could not generate access token")
	}
	refresh, err := crypto.NewRefreshToken(user.uuid, user.name, remember)
	if err != nil {
		return errors.Wrap(err, "could not generate refresh token")
	}

	uid, err := refresh.UserID()
	if err != nil {
		return errors.Wrap(err, "could not get userid from new refresh token")
	}

	refreshModel := storage.NewRefreshTokenModel(uid, refresh.Encrypted, refresh.Expiration, remember)
	err = s.store.CreateRefreshToken(refreshModel)
	if err != nil {
		return errors.Wrap(err, "could not write new refresh token to db")
	}

	tokenC := access.AddToCookie()
	c.SetCookie(tokenC)

	refreshC := refresh.AddToCookie()
	c.SetCookie(refreshC)

	return nil
}

func (s *HandlerSession) handlePostLogout(c echo.Context) error {
	// delete refresh token from db
	refresh, err := c.Cookie("refresh")
	if err == nil && refresh != nil {
		err = s.store.DeleteRefreshTokenByToken(refresh.Value)
		if err != nil {
			fmt.Printf("did not delete refresh token from db: %v\n", err)
		}
	}

	clearCookie("token", "/", c)
	clearCookie("refresh", "/token", c)

	return s.redirectToLogin(c)
}

func clearCookie(name, path string, c echo.Context) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return
	}
	cookie.Value = ""
	cookie.Path = path //very important
	cookie.MaxAge = -1
	c.SetCookie(cookie)
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
