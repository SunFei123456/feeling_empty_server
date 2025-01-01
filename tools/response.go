package tools

import (
  "github.com/labstack/echo/v4"
  "net/http"
)

type Response struct {
  Code int         `json:"code"`
  Msg  string      `json:"msg"`
  Data interface{} `json:"data"`
}

// 成功的响应
func SuccessResponse(c echo.Context, data interface{}) error {
  response := Response{
    Code: http.StatusOK,
    Msg:  "Request Success",
    Data: data,
  }
  return c.JSON(http.StatusOK, echo.Map{"data": response})
}

// 失败的响应
func ErrorResponse(c echo.Context, code int, msg string) error {
  response := Response{
    Code: code,
    Msg:  msg,
    Data: nil,
  }
  return c.JSON(code, echo.Map{"error": response})
}

type PagedResponse struct {
  Code int         `json:"code"`
  Msg  string      `json:"msg"`
  Data interface{} `json:"data"`
  Meta interface{} `json:"meta"`
}

func PagedOkResponse(c echo.Context, data any, meta any) error {
  response := PagedResponse{
    Code: http.StatusOK,
    Msg:  "Request Success",
    Data: data,
    Meta: meta,
  }
  return c.JSON(http.StatusOK, echo.Map{"data": response})
}
