package handle

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
)

type CognitiveExpansionHandler struct {
}

// 1. 获取congnitive_expansion的tags
func (h *CognitiveExpansionHandler) GetTags(c echo.Context) error {
  var tags []string

  if err := db.DB.Table("cognitive_expansion_articles").Select("tags").Distinct().Scan(&tags).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, tags)
}

// 2. 根据pathParme 参数(tag) 来获取指定tag下的所有文章list
func (h *CognitiveExpansionHandler) GetArticlesByTag(c echo.Context) error {
  tag := c.QueryParam("tag")
  var articles []model.CognitiveExpansionArticle
  if err := db.DB.Where("tags =?", tag).Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, articles)
}

// 3. 根据id 获取指定文章
func (h *CognitiveExpansionHandler) GetArticleById(c echo.Context) error {
  id := c.Param("id")
  var article model.CognitiveExpansionArticle
  if err := db.DB.Where("id =?", id).First(&article).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, article)
}

// 4. 获取最近的10篇文章
func (h *CognitiveExpansionHandler) GetLatestArticles(c echo.Context) error {
  var articles []model.CognitiveExpansionArticle
  if err := db.DB.Order("created_at desc").Order("id desc").Limit(10).Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, articles)
}

// 根据title 全文搜素 模糊匹配
func (h *CognitiveExpansionHandler) SearchByTitle(c echo.Context) error {
  title := c.QueryParam("title")
  var articles []model.CognitiveExpansionArticle
  if err := db.DB.Where("title LIKE?", "%"+title+"%").Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, articles)
}
