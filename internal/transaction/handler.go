package transaction

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (handler Handler) Register(e *echo.Echo) {
	e.POST("/households/:id/transactions", handler.Create)
	e.GET("/households/:id/transactions", handler.List)
	e.GET("/households/:id/transactions/:tid", handler.Get)
	e.PATCH("/households/:id/transactions/:tid", handler.Update)
	e.DELETE("/households/:id/transactions/:tid", handler.Delete)
}

func (handler Handler) Create(c *echo.Context) error {
	var input CreateInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "invalid JSON body",
			Code:    "invalid_argument",
		})
	}

	created, err := handler.service.Create(c.Request().Context(), c.Param("id"), input)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusCreated, created)
}

func (handler Handler) List(c *echo.Context) error {
	filter, err := listFilter(c)
	if err != nil {
		return handleError(c, err)
	}

	transactions, err := handler.service.List(c.Request().Context(), c.Param("id"), filter)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, transactions)
}

func (handler Handler) Get(c *echo.Context) error {
	transaction, err := handler.service.Get(c.Request().Context(), c.Param("id"), c.Param("tid"))
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, transaction)
}

func (handler Handler) Update(c *echo.Context) error {
	var input UpdateInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "invalid JSON body",
			Code:    "invalid_argument",
		})
	}

	updated, err := handler.service.Update(c.Request().Context(), c.Param("id"), c.Param("tid"), input)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, updated)
}

func (handler Handler) Delete(c *echo.Context) error {
	if err := handler.service.Delete(c.Request().Context(), c.Param("id"), c.Param("tid")); err != nil {
		return handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func listFilter(c *echo.Context) (ListFilter, error) {
	from, err := optionalTime(c.QueryParam("from"))
	if err != nil {
		return ListFilter{}, ValidationError{Message: "from must be RFC3339"}
	}
	to, err := optionalTime(c.QueryParam("to"))
	if err != nil {
		return ListFilter{}, ValidationError{Message: "to must be RFC3339"}
	}

	return ListFilter{
		From:       from,
		To:         to,
		PayerID:    c.QueryParam("payer_id"),
		CategoryID: c.QueryParam("category_id"),
		Type:       Type(c.QueryParam("type")),
	}, nil
}

func optionalTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func handleError(c *echo.Context, err error) error {
	var validationError ValidationError
	switch {
	case errors.As(err, &validationError):
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: validationError.Message,
			Code:    "invalid_argument",
		})
	case errors.Is(err, ErrNotFound):
		return c.JSON(http.StatusNotFound, errorResponse{
			Message: "transaction not found",
			Code:    "not_found",
		})
	case errors.Is(err, ErrConflict):
		return c.JSON(http.StatusConflict, errorResponse{
			Message: "transaction was already updated",
			Code:    "conflict",
		})
	default:
		slog.Error("transaction request failed", "error", err)
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: "internal server error",
			Code:    "internal",
		})
	}
}
