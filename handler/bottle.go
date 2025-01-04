package handler

import (
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
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

  return OkResponse(c, "创建成功!")
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
    bottleMap := tools.ToMap(bottle, "id", "title", "content", "image_url", "audio_url",
      "mood", "topic_id", "created_at", "resonances", "views", "user")
    bottleMap["user"] = tools.ToMap(bottle.User, "id", "nickname", "avatar", "sex")

    result = append(result, bottleMap)
  }

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

  // 数据清洗

  var result []map[string]any
  for _, bottle := range bottles {
    bottleMap := tools.ToMap(bottle, "id", "title", "content", "image_url", "audio_url",
      "mood", "topic_id", "created_at", "resonances", "views", "user")
    bottleMap["user"] = tools.ToMap(bottle.User, "id", "nickname", "avatar", "sex")
    bottleMap["topic"] = tools.ToMap(bottle.Topic, "id", "title")
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

// HandleGetHotBottles 获取热门漂流瓶 (可选参数 页码, 页数, 时间范围)
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

  // 构建基础查询，注意这里不要使用 Model
  query := h.db.Table("bottles").
    Select("bottles.*, users.*, (bottles.views * 0.4 + bottles.resonances * 0.6) as hotness").
    Joins("LEFT JOIN users ON bottles.user_id = users.id").
    Where("bottles.is_public = ?", true)

  // 根据时间范围筛选
  switch params.TimeRange {
  case "day":
    query = query.Where("bottles.created_at >= ?", time.Now().AddDate(0, 0, -1))
  case "week":
    query = query.Where("bottles.created_at >= ?", time.Now().AddDate(0, 0, -7))
  case "month":
    query = query.Where("bottles.created_at >= ?", time.Now().AddDate(0, -1, 0))
  }

  // 计算总数
  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to count bottles")
  }

  // 按热度排序并分页
  offset := (params.Page - 1) * params.PageSize
  type Result struct {
    model.Bottle
    Hotness float64 `json:"hotness"`
    User    model.User
  }
  var results []Result

  err := query.
    Order("hotness DESC").
    Offset(offset).
    Limit(params.PageSize).
    Scan(&results).Error

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "Failed to get hot bottles")
  }

  // 处理返回数据
  var bottles []map[string]interface{}
  for _, result := range results {
    bottleMap := tools.ToMap(&result.Bottle, "id", "title", "content", "image_url", "audio_url", "mood", "topic_id", "user_id", "created_at", "views", "resonances")
    bottleMap["hotness"] = result.Hotness
    bottleMap["user"] = tools.ToMap(&result.User, "id,nickname,avatar,sex")
    bottles = append(bottles, bottleMap)
  }

  return PagedOkResponse(c, map[string]interface{}{
    "bottles":    bottles,
    "time_range": params.TimeRange,
  }, total, params.Page, params.PageSize)
}
