package handle

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
  "strings"
)

type WebsiteHandler struct{}

// todo 以后数据多了 再加分页, 现在先不考虑了
// 根据category 获取网站卡片列表
func (w *WebsiteHandler) Get(c echo.Context) error {
  // 获取路径参数
  category := c.Param("category")
  var websites []model.Website
  err := db.DB.Where("category = ?", category).Find(&websites)
  if err.Error != nil {
    return c.JSON(http.StatusNotFound, ErrorResponse(c, http.StatusNotFound, err.Error.Error()))
  }
  return SuccessResponse(c, websites)
}

// GetTagsByCategory 获取各个大模块下的标签集合
func (w *WebsiteHandler) GetTagsByCategory(c echo.Context) error {
  category := c.Param("category") // 从查询参数获取类别

  // 查询标签集合
  var tags string
  query := `
        SELECT GROUP_CONCAT(DISTINCT TRIM(tags) SEPARATOR ', ') AS tags
        FROM websites
        WHERE category = ?`

  if err := db.DB.Raw(query, category).Scan(&tags).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  // 去除重复标签
  uniqueTags := uniqueTags(tags)

  // 返回结果
  return SuccessResponse(c, uniqueTags)
}

// 根据标签搜素, 获取包含该tag的网站列表
func (w *WebsiteHandler) SearchByTag(c echo.Context) error {
  // 获取查询参数
  tag := c.QueryParam("tag")

  // 查询包含该标签的网站
  var websites []model.Website
  err := db.DB.Where("tags LIKE?", "%"+tag+"%").Find(&websites).Error
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  return SuccessResponse(c, websites)
}

// 去重函数
func uniqueTags(input string) []string {
  tagSet := make(map[string]struct{})

  for _, tag := range strings.Split(input, ",") {
    trimmedTag := strings.TrimSpace(tag)
    tagSet[trimmedTag] = struct{}{}
  }

  // 将去重后的标签存回切片
  result := make([]string, 0, len(tagSet))
  for tag := range tagSet {
    result = append(result, tag)
  }

  return result
}
