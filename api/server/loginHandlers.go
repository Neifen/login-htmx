package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neifen/htmx-login/api/logging"
	"github.com/neifen/htmx-login/view"
)

type HandlerSession struct {
	store Storage //interfaces are already pointers?
	user  *UserInfo
}

func NewHanderSession(store Storage) *HandlerSession {
	return &HandlerSession{
		store: store,
		user:  EmptyUserInfo(),
	}
}

type UserInfo struct {
	userName string
	uid      string
	crypt    *Crypt
}

func EmptyUserInfo() *UserInfo {
	return &UserInfo{}
}

type Crypt struct {
	token        string
	refreshToken string
}

func (ui *UserInfo) AddCrypt() {
	token, _, err := NewToken(ui)
	if err != nil {
		logging.Error(err)
		return
	}

	refresh := NewRefreshToken(ui)

	ui.crypt = &Crypt{token: token, refreshToken: refresh}
}

func (s *HandlerSession) handleGetLogin(c echo.Context) error {
	if s.isLoggedIn(c.Cookie("token")) {
		return s.redirectToHome(c)
	}

	child := view.Login()
	return view.RenderView(c, child)
}

func (s *HandlerSession) handlePostLogin(c echo.Context) error {
	if s.isLoggedIn(c.Cookie("token")) {
		return s.redirectToHome(c)
	}

	email := c.FormValue("email")
	pw := c.FormValue("password")

	if s.Authenticate(email, pw) {
		//add cookie
		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = s.user.crypt.token
		cookie.Expires = time.Now().Add(15 * time.Minute)
		c.SetCookie(cookie)

		return s.redirectToHome(c)
	}

	// authenticate failed
	return s.redirectToLogin(c)
}

func (s *HandlerSession) handlePostLogout(c echo.Context) error {
	s.user = EmptyUserInfo()

	return s.redirectToLogin(c)
}

func (s *HandlerSession) handleGetRecovery(c echo.Context) error {
	if s.isLoggedIn(c.Cookie("token")) {
		return s.redirectToHome(c)
	}

	child := view.PWRecovery()
	return view.RenderView(c, child)
}

func (s *HandlerSession) handleGetSignup(c echo.Context) error {
	if s.isLoggedIn(c.Cookie("token")) {
		return s.redirectToHome(c)
	}

	child := view.Signup()
	return view.RenderView(c, child)
}

func (s *HandlerSession) handlePostSignup(c echo.Context) error {
	if s.isLoggedIn(c.Cookie("token")) {
		return s.redirectToHome(c)
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
