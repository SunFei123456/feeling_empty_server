package handler

import (
  "fangkong_xinsheng_app/service"
  "fmt"
  "github.com/labstack/echo/v4"
  "net/http"
  "path"
  "time"
)

type QiniuHandler struct{}

// GetUploadToken 获取上传凭证
func (qc *QiniuHandler) GetUploadToken(c echo.Context) error {
  // 生成唯一的文件名
  key := fmt.Sprintf("%d_%s", time.Now().UnixNano(), path.Base(c.QueryParam("filename")))
  token, err := service.GetUploadToken(key)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取上传凭证失败")
  }

  return OkResponse(c, map[string]interface{}{
    "token": token,
    "key":   key,
  })
}
