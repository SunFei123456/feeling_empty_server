package handler

import (
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "strconv"
  "time"
)

type BottleHandler struct {
  db *gorm.DB
}

func NewBottleHandler(db *gorm.DB) *BottleHandler {
  return &BottleHandler{db: db}
}

// HandleCreateBottle 创建漂流瓶
func (h *BottleHandler) HandleCreateBottle(c echo.Context) error {
  var req structs.CreateBottleRequest
  if err := c.Bind(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求体,请检查请求参数"+err.Error())
  }

  // 验证 根据结构体的 validate 标签
  if err := c.Validate(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  // 调用自定义验证方法
  if err := req.Validate(); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  userID := tools.GetUserIDFromContext(c)
  // 创建漂流瓶(uid, content, image_url, audio_url, mood, topic_id, is_public)
  bottle := &model.Bottle{
    Title:    req.Title,
    UserID:   userID,
    Content:  req.Content,
    ImageURL: req.ImageURL,
    AudioURL: req.AudioURL,
    Mood:     req.Mood,
    TopicID:  req.TopicID,
    IsPublic: req.IsPublic,
  }

  if err := h.db.Create(bottle).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "创建漂流瓶失败"+err.Error())
  }

  if err := h.db.Preload("User").First(bottle, "id = ?", bottle.ID).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取漂流瓶失败"+err.Error())
  }

  return OkResponse(c, tools.ToMap(bottle, "id,title,content,image_url,audio_url,mood,topic_id,user_id,created_at,views,user"))
}

// HandleGetRandomBottles 随机获取漂流瓶(10个)
func (h *BottleHandler) HandleGetRandomBottles(c echo.Context) error {
  var bottles []model.Bottle
  if err := h.db.Where("is_public = ?", true).
    Order("RAND()").
    Limit(10).
    Preload("User").
    Find(&bottles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取漂流瓶失败: "+err.Error())
  }

  var result []map[string]any
  for _, bottle := range bottles {
    // 不指定字段列表，将返回所有字段
    // 或者明确指定要返回的字段
    bottleMap := tools.ToMap(bottle, "id", "title", "content", "image_url", "audio_url",
      "mood", "topic_id", "created_at", "views", "user")
    // user 也过滤
    bottleMap["user"] = tools.ToMap(bottle.User, "id", "name", "avatar_url", "sex")

    result = append(result, bottleMap)
  }

  return OkResponse(c, result)
}

// HandleGetBottle 获取漂流瓶详情
func (h *BottleHandler) HandleGetBottle(c echo.Context) error {
  // 将 id 参数转换为 uint
  id, err := strconv.ParseUint(c.Param("id"), 10, 32)
  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的漂流瓶ID")
  }

  var bottle model.Bottle
  if err := h.db.Preload("User").First(&bottle, id).Error; err != nil {
    return ErrorResponse(c, http.StatusNotFound, "漂流瓶不存在"+err.Error())
  }

  userID := tools.GetUserIDFromContext(c)

  // 使用事务来确保数据一致性
  err = h.db.Transaction(func(tx *gorm.DB) error {
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

  var result []map[string]any

  bottleMap := tools.ToMap(bottle, "id", "title", "content", "image_url", "audio_url",
    "mood", "topic_id", "created_at", "views", "user")
  bottleMap["user"] = tools.ToMap(bottle.User, "id", "nickname", "avatar_url", "sex")
  bottleMap["topic"] = tools.ToMap(bottle.Topic, "id", "title", "status")
  result = append(result, bottleMap)

  return OkResponse(c, result)
}

// HandleUpdateBottle 更新漂流瓶
func (h *BottleHandler) HandleUpdateBottle(c echo.Context) error {
  id := c.Param("id")

  var bottle model.Bottle
  if err := h.db.First(&bottle, id).Error; err != nil {
    return ErrorResponse(c, http.StatusNotFound, "漂流瓶不存在"+err.Error())
  }

  // 检查权限
  userID := tools.GetUserIDFromContext(c)
  if bottle.UserID != userID {
    return ErrorResponse(c, http.StatusForbidden, "你没有权限去更新这个漂流瓶")
  }

  var req structs.UpdateBottleRequest
  if err := c.Bind(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数"+err.Error())
  }

  if err := c.Validate(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "参数校验错误"+err.Error())
  }

  updates := make(map[string]any)

  if req.Title != "" {
    updates["title"] = req.Title
  }
  if req.Content != "" {
    updates["content"] = req.Content
  }
  if req.ImageURL != "" {
    updates["image_url"] = req.ImageURL
  }
  if req.AudioURL != "" {
    updates["audio_url"] = req.AudioURL
  }
  if req.Mood != "" {
    updates["mood"] = req.Mood
  }
  if req.IsPublic != nil {
    updates["is_public"] = *req.IsPublic
  }

  if err := h.db.Model(&bottle).Updates(updates).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "更新漂流瓶内容失败"+err.Error())
  }

  return OkResponse(c, "漂流瓶更新成功")
}

// HandleDeleteBottle 删除漂流瓶
func (h *BottleHandler) HandleDeleteBottle(c echo.Context) error {
  id := c.Param("id")
  var bottle model.Bottle
  if err := h.db.Delete(&bottle, id).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "删除漂流瓶失败"+err.Error())
  }
  return OkResponse(c, "漂流瓶删除成功!")
}

// HandleGetBottles 获取漂流瓶列表
func (h *BottleHandler) HandleGetBottles(c echo.Context) error {
  var params structs.BottleQueryParams
  if err := c.Bind(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
  }

  if err := c.Validate(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  if params.Page == 0 {
    params.Page = 1
  }
  if params.PageSize == 0 {
    params.PageSize = 10
  }

  query := h.db.Model(&model.Bottle{}).Preload("User")

  if params.UserID != 0 {
    query = query.Where("user_id = ?", params.UserID)
  }
  if params.TopicID != 0 {
    query = query.Where("topic_id = ?", params.TopicID)
  }
  if params.IsPublic != nil {
    query = query.Where("is_public = ?", *params.IsPublic)
  }

  // 计算总数
  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "统计漂流瓶数量失败"+err.Error())
  }

  // 排序
  if params.Sort != "" {
    query = query.Order(params.Sort)
  } else {
    query = query.Order("created_at DESC")
  }

  // 分页
  offset := (params.Page - 1) * params.PageSize
  var bottles []model.Bottle
  if err := query.Offset(offset).Limit(params.PageSize).Find(&bottles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to get bottles")
  }

  return PagedOkResponse(c, bottles, total, params.Page, params.PageSize)
}

// HandleGetViewedBottles 获取用户查看过的漂流瓶列表
func (h *BottleHandler) HandleGetViewedBottles(c echo.Context) error {
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
  query := h.db.Model(&model.BottleView{}).
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
      "mood", "topic_id", "created_at", "views", "user")
    // user 也过滤
    bottleMap["user"] = tools.ToMap(bottleView.User, "id", "nickname", "avatar_url", "sex")

    result = append(result, bottleMap)
  }

  return PagedOkResponse(c, result, total, params.Page, params.PageSize)
}

// HandleGetRecentViewedBottles 获取最近3天查看的漂流瓶
func (h *BottleHandler) HandleGetRecentViewedBottles(c echo.Context) error {
  userID := tools.GetUserIDFromContext(c)

  // 计算3天前的时间
  threeDaysAgo := time.Now().AddDate(0, 0, -3)

  // 查询最近3天查看过的漂流瓶
  var bottleViews []model.BottleView
  err := h.db.Model(&model.BottleView{}).
    Joins("LEFT JOIN bottles ON bottle_views.bottle_id = bottles.id").
    Where("bottle_views.user_id = ? AND bottle_views.updated_at >= ?", userID, threeDaysAgo).
    Preload("Bottle"). // 预加载漂流瓶信息
    Preload("Bottle.User"). // 预加载漂流瓶作者信息
    Order("bottle_views.updated_at DESC").
    Find(&bottleViews).Error

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取最近查看的漂流瓶失败: "+err.Error())
  }

  // 获取总数
  var total int64
  h.db.Model(&model.BottleView{}).
    Where("user_id = ? AND updated_at >= ?", userID, threeDaysAgo).
    Count(&total)

  var result []map[string]interface{}
  for _, view := range bottleViews {
    if view.Bottle.ID == 0 { // 跳过已删除的漂流瓶
      continue
    }

    bottleMap := tools.ToMap(view.Bottle, "id", "title", "content", "image_url", "audio_url",
      "mood", "topic_id", "created_at", "views")

    // 添加用户信息
    if view.Bottle.User.ID != 0 {
      bottleMap["user"] = tools.ToMap(view.Bottle.User, "id", "nickname", "avatar", "sex")
    }

    // 添加查看时间
    bottleMap["viewed_at"] = view.UpdatedAt

    result = append(result, bottleMap)
  }

  return OkResponse(c, map[string]interface{}{
    "bottles":    result,
    "total":      total,
    "start_time": threeDaysAgo,
    "end_time":   time.Now(),
  })
}

// HandleGetHotBottles 获取热门漂流瓶
func (h *BottleHandler) HandleGetHotBottles(c echo.Context) error {
  var params structs.HotBottleQueryParams
  if err := c.Bind(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
  }

  if err := c.Validate(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  if params.Page == 0 {
    params.Page = 1
  }
  if params.PageSize == 0 {
    params.PageSize = 10
  }

  // 构建基础查询
  query := h.db.Model(&model.Bottle{}).
    Where("is_public = ?", true).
    Preload("User")

  // 根据时间范围筛选
  switch params.TimeRange {
  case "day":
    query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -1))
  case "week":
    query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
  case "month":
    query = query.Where("created_at >= ?", time.Now().AddDate(0, -1, 0))
  }

  // 计算热度分数
  // 热度分数 = 浏览量 * 0.4 + 共鸣值 * 0.6
  query = query.Select("*, (views * 0.4 + resonance_value * 0.6) as hot_score")

  // 计算总数
  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to count bottles")
  }

  // 按热度分数排序并分页
  offset := (params.Page - 1) * params.PageSize
  var bottles []struct {
    model.Bottle
    HotScore float64 `json:"hot_score"`
  }

  err := query.
    Order("hot_score DESC").
    Offset(offset).
    Limit(params.PageSize).
    Find(&bottles).Error

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to get hot bottles")
  }

  return PagedOkResponse(c, map[string]interface{}{
    "bottles":    bottles,
    "time_range": params.TimeRange,
  }, total, params.Page, params.PageSize)
}
