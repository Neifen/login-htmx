package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	LOGIN_PATH    string = "/login"
	SIGNUP_PATH   string = "/signup"
	LOGOUT_PATH   string = "/logout"
	RECOVERY_PATH string = "/recovery"

	HOME_PATH           string = "/"
	HOME_SECONDARY_PATH string = "/home"
)

type APIServer struct {
	apiPath string
	store   Storage
}

func NewAPIHandler(path string, s Storage) *APIServer {
	return &APIServer{apiPath: path, store: s}
}

func (api *APIServer) Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Static("/static", "assets")

	s := NewHanderSession(api.store)

	// login
	e.GET(LOGIN_PATH, s.handleGetLogin)
	e.POST(LOGIN_PATH, s.handlePostLogin)

	e.GET(SIGNUP_PATH, s.handleGetSignup)
	e.POST(SIGNUP_PATH, s.handlePostSignup)

	e.POST(LOGOUT_PATH, s.handlePostLogout)
	e.GET(RECOVERY_PATH, s.handleGetRecovery)

	e.Use()
	// home
	e.GET(HOME_PATH, s.handleGetHome)
	e.GET(HOME_SECONDARY_PATH, s.handleGetHome)

	e.Logger.Fatal(e.Start(api.apiPath))
}
