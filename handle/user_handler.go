package handle

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
  "strconv"
)

// 定义一个UserHandler结构体
type UserHandler struct {
}

// POST 创建一个 user 用户
func (h *UserHandler) Create(c echo.Context) error {
  // 从请求中获取数据
  user := model.User{}

  // 绑定表单数据到 user 结构体
  if err := c.Bind(&user); err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "请求参数错误"})
  }

  // 验证性别字段
  if user.Gender != "male" && user.Gender != "female" && user.Gender != "other" {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "性别参数无效"})
  }

  // 将数据写入数据库
  if err := db.DB.Create(&user).Error; err != nil {
    return c.JSON(http.StatusInternalServerError, echo.Map{"error": "创建用户失败"})
  }

  return c.JSON(http.StatusCreated, echo.Map{"message": "创建成功", "user": user})
}

// 根据id 获取用户
func (h *UserHandler) Get(c echo.Context) error {
  idParam := c.Param("id")
  // 转换为int类型
  id, err := strconv.Atoi(idParam)
  if err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "请求参数错误"})
  }
  var user model.User
  if err := db.DB.First(&user, id).Error; err != nil {
    return c.JSON(http.StatusNotFound, echo.Map{"error": "用户不存在"})
  }
  return c.JSON(http.StatusOK, echo.Map{"user": user})
}

// 获取用户发表的文章的总数
func (h *UserHandler) GetArticlesTotal(c echo.Context) error {
  idParam := c.Param("id")
  // 转换为int类型
  id, err := strconv.Atoi(idParam)
  if err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "请求参数错误"})
  }
  var count int64
  if err := db.DB.Model(&model.Article{}).Where("user_id = ?", id).Count(&count).Error; err != nil {
    return c.JSON(http.StatusNotFound, echo.Map{"error": "用户不存在"})
  }
  return c.JSON(http.StatusOK, echo.Map{"count": count})
}
