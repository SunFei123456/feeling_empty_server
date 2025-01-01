package tools

import (
  "github.com/golang-jwt/jwt"
  "github.com/labstack/echo/v4"
  "time"
)

//GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c echo.Context) uint {
  user := c.Get("user").(*jwt.Token)
  claims := user.Claims.(*jwt.MapClaims)
  id := uint((*claims)["user_id"].(float64))
  return id
}

// GenerateJWTToken 生成JWT token
func GenerateJWTToken(id uint) (string, error) {
  claims := jwt.MapClaims{
    "user_id": id,
    "exp":     time.Now().Add(time.Hour * 72).Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString([]byte("fangkongxinsheng_sf"))
}
