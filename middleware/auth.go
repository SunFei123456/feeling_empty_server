package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"fangkong_xinsheng_app/handler"
)

// JWT 创建 JWT 中间件
func JWT() echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &jwt.MapClaims{},
		SigningKey: []byte("fangkongxinsheng_sf"),
		ErrorHandler: func(err error) error {
			return handler.ErrorResponse(nil, http.StatusUnauthorized, "未授权访问")
		},
	}
	return middleware.JWTWithConfig(config)
}
