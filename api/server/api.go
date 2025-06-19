package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	LOGIN_PATH    string = "/login"
	SIGNUP_PATH   string = "/signup"
	LOGOUT_PATH   string = "/logout"
	RECOVERY_PATH string = "/recovery"

	REFRESH_PATH string = "/token/refresh"

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

	e.POST(REFRESH_PATH, s.handlePostTokenRefresh)

	e.Use(pasetoMiddle())
	// home
	e.GET(HOME_PATH, s.handleGetHome)
	e.GET(HOME_SECONDARY_PATH, s.handleGetHome)

	e.Logger.Fatal(e.Start(api.apiPath))
}

func pasetoMiddle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
	}
}
