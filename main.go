package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := newRouter()
	e.Logger.Fatal(e.Start(":8080"))
}

func newRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", postEnter)

	return e
}

func postEnter(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}
