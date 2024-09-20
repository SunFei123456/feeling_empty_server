package handle

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
)

type StatsHandler struct{}

// 获取文章的总浏览量 + 文章总数 + 文章类型数量
func (h *StatsHandler) GetArticleStats(c echo.Context) error {
  var articleTotal int64
  if err := db.DB.Model(&model.Article{}).Count(&articleTotal).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  var viewsTotal int64
  if err := db.DB.Model(&model.Article{}).Select("sum(views) as total_views").Scan(&viewsTotal).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  var categoryTotal int64
  if err := db.DB.Model(&model.Article{}).Select("count(distinct category) as total_category").Scan(&categoryTotal).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  var stats struct {
    ArticleTotal  int64 `json:"article_total"`
    ViewsTotal    int64 `json:"views_total"`
    CategoryTotal int64 `json:"category_total"`
  }
  stats.ArticleTotal = articleTotal
  stats.ViewsTotal = viewsTotal
  stats.CategoryTotal = categoryTotal

  return SuccessResponse(c, stats)
}
