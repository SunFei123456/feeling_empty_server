package handle

import (
  "encoding/json"
  "fmt"
  "github.com/labstack/echo/v4"
  "io"
  "math/rand"
  "net/http"
  "os"
  "path/filepath"
  "time"
)

type ImageHandler struct {
}

// UploadImage 处理图片上传并保存到服务器
func (h *ImageHandler) Upload(c echo.Context) error {
  // 获取上传的文件
  file, err := c.FormFile("image")
  if err != nil {
    return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid file"})
  }

  // 打开上传的文件
  src, err := file.Open()
  if err != nil {
    return err
  }
  defer src.Close()

  // 确保目标文件夹存在
  destFolder := "static/avatars"
  if _, err := os.Stat(destFolder); os.IsNotExist(err) {
    err = os.MkdirAll(destFolder, os.ModePerm)
    if err != nil {
      return err
    }
  }

  // 生成新的文件名（使用时间戳）
  newFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
  filePath := filepath.Join(destFolder, newFileName)

  // 创建目标文件
  dest, err := os.Create(filePath)
  if err != nil {
    return err
  }
  defer dest.Close()

  // 将上传的文件内容复制到目标文件
  if _, err := io.Copy(dest, src); err != nil {
    return err
  }

  // 返回文件路径给前端
  fileURL := "http://localhost:8080/static/avatars/" + newFileName
  fmt.Println("fileURL: ", fileURL)
  return SuccessResponse(c, fileURL)
}

// 获取 Bing 每日图片
// 定义结构体以匹配 Bing API 响应
// BingImageResponse 表示 Bing API 返回的图片信息
type BingImageResponse struct {
  Images []struct {
    // Images 表示 Bing API 返回的每个壁纸的信息
    Url           string `json:"url"`           // 壁纸的相对 URL，用于获取图像
    Title         string `json:"title"`         // 壁纸的标题，描述性文本
    Copyright     string `json:"copyright"`     // 版权信息，指明版权所有者或来源
    StartDate     string `json:"startdate"`     // 壁纸开始使用的日期
    FullStartDate string `json:"fullstartdate"` // 壁纸开始使用的完整日期和时间
    EndDate       string `json:"enddate"`       // 壁纸停止使用的日期
    UrlBase       string `json:"urlbase"`       // 壁纸的基础 URL，通常用于构造完整的图像链接
  } `json:"images"`
}

// GetRandomBingImage 获取 Bing 随机图片
func (h *ImageHandler) GetRandomBingImage(c echo.Context) error {
  // 请求 Bing API，获取 10 张图片
  resp, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=10&mkt=zh-CN")
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  defer resp.Body.Close()

  // 检查响应状态码
  if resp.StatusCode != http.StatusOK {
    return ErrorResponse(c, http.StatusInternalServerError, "请求 Bing API 失败")
  }

  // 解析 JSON 响应
  var bingResponse BingImageResponse
  if err := json.NewDecoder(resp.Body).Decode(&bingResponse); err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  // 随机选择一张图片
  if len(bingResponse.Images) > 0 {
    rand.Seed(time.Now().UnixNano()) // 确保每次都能得到不同的随机数
    randomIndex := rand.Intn(len(bingResponse.Images))
    imageUrl := "https://cn.bing.com" + bingResponse.Images[randomIndex].Url
    res := map[string]string{
      "imageUrl":  imageUrl,
      "title":     bingResponse.Images[randomIndex].Title,
      "copyright": bingResponse.Images[randomIndex].Copyright,
    }
    return SuccessResponse(c, res)
  }

  return ErrorResponse(c, http.StatusInternalServerError, "未找到 Bing 随机图片")
}
