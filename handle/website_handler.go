package handle

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
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
