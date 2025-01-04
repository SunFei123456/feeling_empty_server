package handler

import (
  "fangkong_xinsheng_app/db"
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "strconv"
  "time"
)

type BottleViewHandler struct{}

// HandleCreateBottleView 创建漂流瓶查看记录 (同步更新该漂流瓶的浏览量)
func (bv *BottleViewHandler) HandleCreateBottleView(c echo.Context) error {
  // 将 id 参数转换为 uint
  id, err := strconv.ParseUint(c.QueryParam("id"), 10, 32)
  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的漂流瓶ID")
  }

  var bottle model.Bottle
  if err := db.DB.Preload("User").First(&bottle, id).Error; err != nil {
    return ErrorResponse(c, http.StatusNotFound, "漂流瓶不存在"+err.Error())
  }

  userID := tools.GetUserIDFromContext(c)

  // 使用事务来确保数据一致性
  err = db.DB.Transaction(func(tx *gorm.DB) error {
    // 1. 增加浏览量
    if err := tx.Model(&bottle).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error; err != nil {
      return err
    }

    // 2. 记录查看历史
    bottleView := &model.BottleView{
      BottleID: uint(id),
      UserID:   userID,
    }

    // 如果已经查看过，就更新时间，否则创建新记录
    result := tx.Where("bottle_id = ? AND user_id = ?", id, userID).
      FirstOrCreate(bottleView)

    if result.Error != nil {
      return result.Error
    }

    // 如果记录已存在（RowsAffected = 0），则更新时间戳
    if result.RowsAffected == 0 {
      if err := tx.Model(bottleView).Update("updated_at", time.Now()).Error; err != nil {
        return err
      }
    }
    return nil
  })

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "创建漂流瓶浏览历史失败"+err.Error())
  }

  return OkResponse(c, "创建漂流瓶浏览历史成功!")
}

// HandleDeleteBottleView 删除指定的漂流瓶查看记录
func (bv *BottleViewHandler) HandleDeleteBottleView(c echo.Context) error {
  id := c.Param("id")
  var bottleView model.BottleView
  if err := db.DB.Delete(&bottleView, id).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "删除漂流瓶浏览历史失败"+err.Error())
  }
  return OkResponse(c, "漂流瓶浏览历史删除成功!")
}

// HandleDeleteAllBottleViews 删除用户的全部的漂流瓶浏览历史
func (bv *BottleViewHandler) HandleDeleteAllBottleViews(c echo.Context) error {
  userID := tools.GetUserIDFromContext(c)
  if err := db.DB.Where("user_id = ?", userID).Delete(&model.BottleView{}).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "删除漂流瓶浏览历史失败"+err.Error())
  }
  return OkResponse(c, "漂流瓶浏览历史删除成功!")
}

// HandleGetBottleViews 获取漂流瓶浏览历史
func (bv *BottleViewHandler) HandleGetBottleViews(c echo.Context) error {
  userID := tools.GetUserIDFromContext(c)

  var params structs.BottleQueryParams
  if err := c.Bind(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数"+err.Error())
  }

  if params.Page == 0 {
    params.Page = 1
  }
  if params.PageSize == 0 {
    params.PageSize = 10
  }

  // 查询用户查看过的漂流瓶
  query := db.DB.Model(&model.BottleView{}).
    Select("bottle_views.*, bottles.*").
    Joins("LEFT JOIN bottles ON bottle_views.bottle_id = bottles.id").
    Where("bottle_views.user_id = ?", userID).
    Preload("Bottle.User")

  // 计算总数
  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to count bottles")
  }

  // 排序
  query = query.Order("bottle_views.updated_at DESC")

  // 分页
  offset := (params.Page - 1) * params.PageSize
  var bottleViews []model.BottleView
  if err := query.Offset(offset).Limit(params.PageSize).Find(&bottleViews).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to get bottles")
  }

  var result []map[string]any
  for _, bottleView := range bottleViews {
    // 不指定字段列表，将返回所有字段
    // 或者明确指定要返回的字段
    bottleMap := tools.ToMap(bottleView.Bottle, "id", "title", "content", "image_url", "audio_url",
      "mood", "topic_id", "created_at", "views", "resonances", "user")
    // user 也过滤
    bottleMap["user"] = tools.ToMap(bottleView.Bottle.User, "id", "nickname", "avatar", "sex")

    result = append(result, bottleMap)
  }

  return PagedOkResponse(c, result, total, params.Page, params.PageSize)
}
