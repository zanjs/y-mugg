package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

// ResourceController is
type ResourceController interface {
	GetAll(c *echo.Context)
	Get(c *echo.Context)
	Create(c *echo.Context)
	Update(c *echo.Context)
	Delete(c *echo.Context)
}

type httpError struct {
	code    int
	Key     string `json:"error"`
	Message string `json:"message"`
}

// UnityJSON is Return to unity JSON Data
type UnityJSON struct {
	Status  string      `json:"status"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Code    string      `json:"code"`
}

func newHTTPError(code int, key string, msg string) *httpError {
	return &httpError{
		code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *httpError) Error() string {
	return e.Key + ": " + e.Message
}

// httpErrorHandler customize echo's HTTP error handler.
func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "ServerError"
		msg  string
	)

	if he, ok := err.(*httpError); ok {
		code = he.code
		key = he.Key
		msg = he.Message
	} else if ee, ok := err.(*echo.HTTPError); ok {
		code = ee.Code
		key = http.StatusText(code)
		msg = key
		// } else if config.Debug {
		// 	msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			err := c.JSON(code, newHTTPError(code, key, msg))
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

// Controller is base controller
type Controller struct{}

func (ctl Controller) buildSuccessData(data interface{}) map[string]interface{} {
	return echo.Map{
		"status": "success",
		"error":  false,
		"data":   data,
	}
}

func (ctl Controller) buildErrorData(message string) map[string]interface{} {
	return echo.Map{
		"status":  "error",
		"error":   true,
		"message": message,
	}
}

// SuccessResponse is
func (ctl Controller) SuccessResponse(c echo.Context, data interface{}) {
	c.JSON(200, ctl.buildSuccessData(data))
}

// ErrorResponse is
func (ctl Controller) ErrorResponse(c echo.Context, statusCode int, message string) {
	c.JSON(statusCode, ctl.buildErrorData(message))
}
