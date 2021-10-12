package web

import (
	"net/http"

	"github.com/arata-nvm/monban/domain"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", postEnter)

	return e
}

func postEnter(ctx echo.Context) error {
	var req enterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if req.StudentID == 0 {
		return ctx.JSON(http.StatusBadRequest, enterError{
			Message: "field `student_id` is required.",
		})
	}

	if err := domain.Entry(req.StudentID); err != nil {
		ctx.Logger().Warn(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusNoContent)
}

type enterRequest struct {
	StudentID int `json:"student_id" validate:"required"`
}

type enterError struct {
	Message string `json:"message"`
}
