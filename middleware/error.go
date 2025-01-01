package middleware

import (
  "fangkong_xinsheng_app/handler"
  "github.com/labstack/echo/v4"
  "net/http"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
  code := http.StatusInternalServerError
  message := "Internal Server Error"

  if he, ok := err.(*echo.HTTPError); ok {
    code = he.Code
    message = he.Message.(string)
  }

  _ = handler.ErrorResponse(c, code, message)
}
