package handler

import (
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/service"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "strconv"
)

type TopicHandler struct {
  db                 *gorm.DB
  interactionService *service.BottleInteractionService
}

func NewTopicHandler(db *gorm.DB) *TopicHandler {
  return &TopicHandler{
    db:                 db,
    interactionService: service.NewBottleInteractionService(db),
  }
}

// HandleGetSystemTopics 获取系统话题列表(前9个)
func (h *TopicHandler) HandleGetSystemTopics(c echo.Context) error {
  var topics []model.Topic
  if err := h.db.Where("type = ?", 0).
    Select("id, title").
    Order("created_at DESC").
    Limit(9).
    Find(&topics).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取系统话题失败")
  }

  // 处理返回数据
  var result []map[string]interface{}
  for _, topic := range topics {
    result = append(result, map[string]interface{}{
      "id":    topic.ID,
      "title": topic.Title,
    })
  }

  return OkResponse(c, result)
}

// HandleGetTopicBottles 获取话题下的漂流瓶列表
func (h *TopicHandler) HandleGetTopicBottles(c echo.Context) error {
  var params structs.TopicQueryParams
  if err := c.Bind(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
  }

  // 设置默认分页参数
  if params.Page <= 0 {
    params.Page = 1
  }
  if params.PageSize <= 0 {
    params.PageSize = 10
  }

  topicID := c.Param("id")
  sortBy := c.QueryParam("sort") // "new" 或 "hot"
  userID := tools.GetUserIDFromContext(c)

  query := h.db.Model(&model.BottleTopic{}).
    Joins("LEFT JOIN bottles ON bottle_topics.bottle_id = bottles.id").
    Where("bottle_topics.topic_id = ?", topicID).
    Preload("Bottle").
    Preload("Bottle.User")

  // 计算总数
  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取漂流瓶数量失败")
  }

  // 根据排序方式设置排序
  if sortBy == "hot" {
    query = query.Order("bottles.views * 0.4 + bottles.resonances * 0.6 DESC")
  } else {
    query = query.Order("bottle_topics.created_at DESC")
  }

  // 分页
  var bottleTopics []model.BottleTopic
  if err := query.Offset((params.Page - 1) * params.PageSize).
    Limit(params.PageSize).
    Find(&bottleTopics).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取漂流瓶列表失败")
  }

  // 处理返回数据
  var result []map[string]interface{}
  for _, bt := range bottleTopics {
    bottleMap := tools.ToMap(&bt.Bottle, "id", "title", "content", "image_url", "audio_url",
      "mood", "topic_id", "created_at", "views", "resonances", "shares", "favorites")

    // 使用服务添加交互状态
    h.interactionService.EnrichBottleWithInteractionStatus(bottleMap, userID, bt.BottleID)

    if bt.Bottle.User.ID != 0 {
      bottleMap["user"] = tools.ToMap(&bt.Bottle.User, "id", "nickname", "avatar", "sex")
    }
    result = append(result, bottleMap)
  }

  return PagedOkResponse(c, result, total, params.Page, params.PageSize)
}

// HandleGetTopicInfo 获取话题详细信息
func (h *TopicHandler) HandleGetTopicInfo(c echo.Context) error {
  topicID := c.Param("id")

  var topic model.Topic
  if err := h.db.First(&topic, topicID).Error; err != nil {
    return ErrorResponse(c, http.StatusNotFound, "话题不存在")
  }

  // 获取内容数量
  var contentCount int64
  h.db.Model(&model.BottleTopic{}).Where("topic_id = ?", topicID).Count(&contentCount)

  // 获取参与人数
  var participantCount int64
  h.db.Model(&model.BottleTopic{}).
    Where("topic_id = ?", topicID).
    Distinct("user_id").
    Count(&participantCount)

  return OkResponse(c, map[string]interface{}{
    "id":                topic.ID,
    "title":             topic.Title,
    "desc":              topic.Desc,
    "views":             topic.Views,
    "content_count":     contentCount,
    "participant_count": participantCount,
  })
}

// HandleGetHotTopics 获取内容最多的前n个话题
func (h *TopicHandler) HandleGetHotTopics(c echo.Context) error {
  // 定义请求参数 limit n
  limit := c.QueryParam("limit")
  // 转换为整数
  limitInt, err := strconv.Atoi(limit)
  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
  }

  type Result struct {
    ID           uint   `json:"id"`
    Title        string `json:"title"`
    ContentCount int64  `json:"content_count"`
    BgImage      string `json:"bg_image"`
    View         int    `json:"views"`
  }

  var results []Result
  err = h.db.Model(&model.BottleTopic{}).
    Select("topic_id as id, topics.title, count(*) as content_count, topics.bg_image, topics.views").
    Joins("LEFT JOIN topics ON bottle_topics.topic_id = topics.id").
    Group("topic_id").
    Order("content_count DESC").
    Limit(limitInt).
    Scan(&results).Error

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取热门话题失败")
  }

  return OkResponse(c, results)
}

// HandleCreateTopic 创建话题
func (h *TopicHandler) HandleCreateTopic(c echo.Context) error {
  var req structs.CreateTopicRequest
  if err := c.Bind(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
  }

  if err := c.Validate(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  if req.BgImage == "" {
    req.BgImage = "https://www.zmtc.com/wp-content/uploads/2023/0309/20230309070220562.jpg"
  }

  topic := &model.Topic{
    Title:   req.Title,
    Type:    req.Type,
    BgImage: req.BgImage,
  }

  if err := h.db.Create(topic).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "创建话题失败")
  }

  return OkResponse(c, topic)
}

// HandleGetAllTopics 获取所有话题
func (h *TopicHandler) HandleGetAllTopics(c echo.Context) error {
  // 获取搜索关键词
  keyword := c.QueryParam("keyword")

  // 构建查询
  query := h.db.Select("id, title").Order("created_at DESC")

  // 如果有关键词，添加模糊查询条件
  if keyword != "" {
    query = query.Where("title LIKE ?", "%"+keyword+"%")
  }

  var topics []model.Topic
  if err := query.Find(&topics).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取话题列表失败")
  }

  // 处理返回数据
  var result []map[string]interface{}
  for _, topic := range topics {
    result = append(result, map[string]interface{}{
      "id":    topic.ID,
      "title": topic.Title,
    })
  }

  return OkResponse(c, result)
}

// HandleSearchTopics 搜索话题
func (h *TopicHandler) HandleSearchTopics(c echo.Context) error {
  // 获取搜索关键词
  keyword := c.QueryParam("keyword")
  if keyword == "" {
    return ErrorResponse(c, http.StatusBadRequest, "搜索关键词不能为空")
  }

  var topics []model.Topic
  if err := h.db.Select("id,title,views,bg_image").
    Where("title LIKE ?", "%"+keyword+"%").
    Order("views DESC").
    Find(&topics).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "搜索话题失败")
  }

  // 处理返回数据
  var result []map[string]interface{}
  for _, topic := range topics {
    result = append(result, map[string]interface{}{
      "id":       topic.ID,
      "title":    topic.Title,
      "views":    topic.Views,
      "bg_image": topic.BgImage,
    })
  }

  return OkResponse(c, result)
}

// HandleAddTopicViews 增加话题浏览量
func (h *TopicHandler) HandleAddTopicViews(c echo.Context) error {
  topicID := c.Param("id")
  err := h.db.Model(&model.Topic{}).Where("id = ?", topicID).
    Update("views", gorm.Expr("views + ?", 1)).Error
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "增加话题浏览量失败")
  }
  return OkResponse(c, "增加话题浏览量成功")
}
