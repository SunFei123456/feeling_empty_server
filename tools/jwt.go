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

// GenerateJWTToken 生成JWT 返回token和过期时间戳
func GenerateJWTToken(id uint) (string, int64, error) {
  expirationTime := time.Now().Add(time.Hour * 72)
  claims := jwt.MapClaims{
    "user_id": id,
    "exp":     expirationTime.Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signedToken, err := token.SignedString([]byte("fangkongxinsheng_sf"))
  if err != nil {
    return "", 0, err // 返回空字符串和零时间戳，同时返回错误
  }
  return signedToken, expirationTime.Unix(), nil // 返回签名的 token 和过期时间戳
}
