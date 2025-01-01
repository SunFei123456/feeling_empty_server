package handler

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

// Response 统一响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// PagedResponse 分页响应结构
type PagedResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
    Total   int64      `json:"total"`
    Page    int        `json:"page"`
    Size    int        `json:"size"`
}

// OkResponse 成功响应
func OkResponse(c echo.Context, data interface{}) error {
    return c.JSON(http.StatusOK, Response{
        Code:    http.StatusOK,
        Message: "success",
        Data:    data,
    })
}

// ErrorResponse 错误响应
func ErrorResponse(c echo.Context, status int, message string) error {
    return c.JSON(status, Response{
        Code:    status,
        Message: message,
    })
}

// PagedOkResponse 分页成功响应
func PagedOkResponse(c echo.Context, data interface{}, total int64, page, size int) error {
    return c.JSON(http.StatusOK, PagedResponse{
        Code:    http.StatusOK,
        Message: "success",
        Data:    data,
        Total:   total,
        Page:    page,
        Size:    size,
    })
} 