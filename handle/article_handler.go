package handle

import (
  "encoding/json"
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
  "quick-start/structs"
  "quick-start/tools"
  "strconv"
)

type ArticleHandler struct{}

// 创建文章
func (h *ArticleHandler) Create(c echo.Context) error {
  // 获取参数
  title := c.FormValue("title")
  content := c.FormValue("content")
  category := c.FormValue("category")
  tags := c.FormValue("tags")
  coverImage := tools.GetCoverImage(category)

  // 将 tags 转换为 JSON 字符串
  tagsJSON, err := json.Marshal([]string{tags})
  if err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid tags format"})
  }
  // 创建文章
  article := model.Article{
    Title:      title,
    Content:    content,
    CoverImage: coverImage,
    Category:   category,
    Tags:       string(tagsJSON),
    UserID:     1,
  }

  // 绑定表单数据到 user 结构体
  if err := c.Bind(&article); err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "请求参数错误"})
  }

  // 将数据写入数据库
  if err := db.DB.Create(&article).Error; err != nil {
    return c.JSON(http.StatusInternalServerError, echo.Map{"error": "创建文章失败"})
  }

  return c.JSON(http.StatusCreated, echo.Map{"message": "创建成功", "article": article})
}

// 根据类别获取文章(limit 5)
func (h *ArticleHandler) Get(c echo.Context) error {
  category := c.Param("category")

  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Where("category = ?", category).Limit(5).Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, articles)
}

// 根据id获取文章
func (h *ArticleHandler) GetOne(c echo.Context) error {
  idParam := c.Param("id")
  // 转换为int类型
  id, err := strconv.Atoi(idParam)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  var article model.Article
  if err := db.DB.Preload("UserInfo").First(&article, id).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  data := tools.ToMap(article, "id", "title", "content", "coverImage", "tags", "category", "views", "created_at", "updated_at")
  data["user"] = tools.ToMap(article.UserInfo, "id", "nickname", "avatar", "address")
  return SuccessResponse(c, data)

}

// 获取最近的文章
func (h *ArticleHandler) GetLatest(c echo.Context) error {
  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Order("created_at desc").Limit(5).Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, articles)
}

// 根据类别获取文章列表
func (h *ArticleHandler) GetArticlesListByCategory(c echo.Context) error {
  category := c.Param("category")
  if category == "Recently" {
    return h.GetLatest(c)
  }
  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Where("category = ?", category).Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, articles)
}

// 随机在热门帖子中获取3篇文章
func (h *ArticleHandler) GetRandomHotArticles(c echo.Context) error {
  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Where("is_hot = ?", true).Order("RAND()").Limit(3).Find(&articles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  var data []map[string]interface{}
  for _, article := range articles {
    item := tools.ToMap(article, "id", "title", "content", "coverImage", "tags", "category", "created_at", "updated_at")
    item["user"] = tools.ToMap(article.UserInfo, "id", "nickname", "avatar", "address")
    data = append(data, item)
  }
  return SuccessResponse(c, data)
}

// 增加文章的浏览量
func (h *ArticleHandler) IncrementViewCount(c echo.Context) error {
  articleID := c.Param("id")
  var article model.Article

  // Find the article by ID
  if err := db.DB.First(&article, articleID).Error; err != nil {
    return ErrorResponse(c, http.StatusNotFound, "Article not found")
  }

  // Increment the view count
  article.Views++

  // Save the updated article
  if err := db.DB.Save(&article).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  return SuccessResponse(c, article.Views)
}

// 根据文章的id获取评论列表(分页)
func (h *ArticleHandler) GetCommentsByArticleID(c echo.Context) error {
  articleID := c.Param("id")
  // 获取分页参数
  page, _ := strconv.Atoi(c.QueryParam("page")) // 页码
  size := 7                                     // 每页大小
  // 计算当前页的起始位置
  offset := (page - 1) * size // 计算偏移量

  // 如果没有传页码，默认第一页
  var comments []model.Comment
  if err := db.DB.Where("commentable_id = ? and commentable_type = 'article' and status = 'visible'", articleID).
    Order("created_at DESC").Limit(size).Offset(offset).Find(&comments).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  // 组装分页器数据
  var total int64
  if err := db.DB.Model(&model.Comment{}).Where("commentable_id = ? and commentable_type = 'article' and status = 'visible'", articleID).Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  // 计算分页信息
  meta := structs.NewPagination(total, page, size)

  return PagedOkResponse(c, comments, meta)

}

// 获取发表文章的总数
func (h *ArticleHandler) GetArticlesTotal(c echo.Context) error {
  userID := c.Param("id")
  var count int64
  if err := db.DB.Model(&model.Article{}).Where("user_id =?", userID).Count(&count).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, count)
}
