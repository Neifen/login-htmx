package view

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func RenderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	if c.Request().Header.Get("HX-Request") != "true" {
		// whole page
		return BuildBase(nil, nil, cmp).Render(c.Request().Context(), c.Response().Writer)
	}

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func ReplaceUrl(path string, c echo.Context, cmp templ.Component) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		// standard redirect
		return c.Redirect(http.StatusSeeOther, path)
	}

	c.Response().Header().Set("HX-Replace-Url", path)
	return RenderView(c, cmp)
}