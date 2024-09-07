package handle

import (
  "encoding/json"
  "fmt"
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
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
  fmt.Println("category:", category)

  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Where("category = ?", category).Limit(4).Find(&articles).Error; err != nil {
    return c.JSON(http.StatusNotFound, echo.Map{"error": "文章不存在"})
  }
  return c.JSON(http.StatusOK, echo.Map{"data": articles})
}

// 根据id获取文章
func (h *ArticleHandler) GetOne(c echo.Context) error {
  idParam := c.Param("id")
  // 转换为int类型
  id, err := strconv.Atoi(idParam)
  if err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "请求参数错误"})
  }
  var article model.Article
  if err := db.DB.Preload("UserInfo").First(&article, id).Error; err != nil {
    return c.JSON(http.StatusNotFound, echo.Map{"error": "文章不存在"})
  }

  data := tools.ToMap(article, "id", "title", "content", "coverImage", "tags", "category", "created_at", "updated_at")
  data["user"] = tools.ToMap(article.UserInfo, "id", "nickname", "avatar", "address")
  return c.JSON(http.StatusOK, echo.Map{"data": data})

}

// 获取最近的文章
func (h *ArticleHandler) GetLatest(c echo.Context) error {
  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Order("created_at desc").Limit(5).Find(&articles).Error; err != nil {
    return c.JSON(http.StatusNotFound, echo.Map{"error": "文章不存在"})
  }
  return c.JSON(http.StatusOK, echo.Map{"data": articles})
}

// 根据类别获取文章列表
func (h *ArticleHandler) GetArticlesListByCategory(c echo.Context) error {
  category := c.Param("category")
  fmt.Println("category:", category)
  if category == "Recently" {
    return h.GetLatest(c)
  }
  var articles []model.Article
  if err := db.DB.Preload("UserInfo").Where("category = ?", category).Find(&articles).Error; err != nil {
    return c.JSON(http.StatusNotFound, echo.Map{"error": "文章不存在"})
  }
  return c.JSON(http.StatusOK, echo.Map{"data": articles})
}
