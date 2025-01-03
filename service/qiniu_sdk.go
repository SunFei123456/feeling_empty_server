package service

import (
  "fmt"
  "github.com/qiniu/go-sdk/v7/auth/qbox"
  "github.com/qiniu/go-sdk/v7/storage"
  _ "time"
)

type QiniuConfig struct {
  AccessKey    string
  SecretKey    string
  BucketName   string
  BucketDomain string
}

var defaultConfig *QiniuConfig

// InitQiniu 初始化七牛云配置
func InitQiniu(cfg *QiniuConfig) {
  defaultConfig = cfg
}

// GetUploadToken 获取上传凭证
func GetUploadToken(key string) (string, error) {
  if defaultConfig == nil {
    return "", fmt.Errorf("七牛云配置未初始化")
  }

  mac := qbox.NewMac(defaultConfig.AccessKey, defaultConfig.SecretKey)

  putPolicy := storage.PutPolicy{
    Scope:      fmt.Sprintf("%s:%s", defaultConfig.BucketName, key),
    ReturnBody: `{"key":"$(key)","hash":"$(etag)","size":$(fsize),"name":"$(fname)"}`,
    Expires:    3600, // 1小时有效期
  }
  return putPolicy.UploadToken(mac), nil
}

// GetResourceURL 生成资源的完整URL
func GetResourceURL(key string) string {
  return fmt.Sprintf("https://%s/%s", defaultConfig.BucketDomain, key)
}
