package handler

import (
  "context"
  "fmt"
  "github.com/labstack/echo/v4"
  "github.com/tencentyun/cos-go-sdk-v5"
  "net/http"
  "net/url"
  "os"
  "strings"
  "time"
)

type COSHandler struct{}

// 允许上传的文件类型及其对应的Content-Type
var allowedExtensions = map[string]string{
  "jpg":  "image/jpeg",
  "jpeg": "image/jpeg",
  "png":  "image/png",
  "gif":  "image/gif",
  "mp3":  "audio/mpeg",
  "wav":  "audio/wav",
  "m4a":  "audio/mp4",
}

// GetUploadToken 获取上传预签名URL
func (h *COSHandler) GetUploadToken(c echo.Context) error {
  // 获取并验证文件扩展名
  ext := strings.ToLower(strings.TrimPrefix(c.QueryParam("ext"), "."))
  if ext == "" {
    return ErrorResponse(c, http.StatusBadRequest, "文件扩展名不能为空")
  }

  contentType, ok := allowedExtensions[ext]
  if !ok {
    return ErrorResponse(c, http.StatusBadRequest, "不支持的文件类型，仅支持jpg/jpeg/png/gif/mp3/wav/m4a格式")
  }

  // 生成唯一的文件名，按类型分目录
  var folder string
  switch ext {
  case "jpg", "jpeg", "png", "gif":
    folder = "images"
  case "mp3", "wav", "m4a":
    folder = "audios"
  }

  key := fmt.Sprintf("%s/%d.%s", folder, time.Now().UnixNano(), ext)

  // 构建 COS 客户端
  bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", os.Getenv("COS_BUCKET"), os.Getenv("COS_REGION"))
  u, _ := url.Parse(bucketURL)
  b := &cos.BaseURL{BucketURL: u}

  client := cos.NewClient(b, &http.Client{
    Transport: &cos.AuthorizationTransport{
      SecretID:  os.Getenv("COS_SECRET_ID"),
      SecretKey: os.Getenv("COS_SECRET_KEY"),
    },
  })

  // 构建预签名URL的选项
  opt := &cos.PresignedURLOptions{
    Header: &http.Header{},
    Query:  &url.Values{},
  }
  opt.Header.Add("Content-Type", contentType)

  // 获取预签名URL
  presignedURL, err := client.Object.GetPresignedURL(
    context.Background(),
    http.MethodPut,
    key,
    os.Getenv("COS_SECRET_ID"),
    os.Getenv("COS_SECRET_KEY"),
    time.Hour,
    opt,
  )
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "生成预签名URL失败: "+err.Error())
  }

  return OkResponse(c, map[string]interface{}{
    "url":          presignedURL.String(),
    "key":          key,
    "bucket":       os.Getenv("COS_BUCKET"),
    "region":       os.Getenv("COS_REGION"),
    "content_type": contentType,
    "expired_time": time.Now().Add(time.Hour).Unix(),
  })
}
