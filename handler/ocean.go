package handler

import (
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "fangkong_xinsheng_app/service"
)

type OceanHandler struct {
  db *gorm.DB
  interactionService *service.BottleInteractionService
}

func NewOceanHandler(db *gorm.DB) *OceanHandler {
  return &OceanHandler{
    db: db,
    interactionService: service.NewBottleInteractionService(db),
  }
}

// HandleGetOceanBottles 获取指定海域下的随机50个瓶子信息
func (h *OceanHandler) HandleGetOceanBottles(c echo.Context) error {
  userID := tools.GetUserIDFromContext(c)
  // 从路径参数获取海域ID
  oceanID := c.Param("ocean_id")

  // 构建查询
  query := h.db.Model(&model.OceanBottle{}).
    Joins("LEFT JOIN bottles ON ocean_bottles.bottle_id = bottles.id").
    Where("ocean_bottles.ocean_id = ?", oceanID).
    Preload("Bottle").
    Preload("Bottle.User")

  // 获取随机50个瓶子
  var oceanBottles []model.OceanBottle
  if err := query.
    Order("RAND()").
    Limit(50).
    Find(&oceanBottles).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取瓶子列表失败")
  }

  // 处理返回数据
  var result []map[string]interface{}
  for _, ob := range oceanBottles {
    bottleMap := tools.ToMap(&ob.Bottle, "id", "title", "content", "image_url", "audio_url", "mood", "topic_id", "created_at", "views", "resonances", "shares", "favorites")

    // 使用服务添加交互状态
    h.interactionService.EnrichBottleWithInteractionStatus(bottleMap, userID, ob.BottleID)

    if ob.Bottle.User.ID != 0 {
      bottleMap["user"] = tools.ToMap(&ob.Bottle.User, "id", "nickname", "avatar", "sex")
    }
    result = append(result, bottleMap)
  }

  return OkResponse(c, result)
}

// HandleGetOceans 获取所有海域信息及瓶子数量
func (h *OceanHandler) HandleGetOceans(c echo.Context) error {
  var oceans []model.Ocean

  // 获取所有海域
  if err := h.db.Find(&oceans).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取海域列表失败")
  }

  // 获取每个海域的瓶子数量
  for i := range oceans {
    var count int64
    if err := h.db.Model(&model.OceanBottle{}).
      Where("ocean_id = ?", oceans[i].ID).
      Count(&count).Error; err != nil {
      return ErrorResponse(c, http.StatusInternalServerError, "获取瓶子数量失败")
    }
    oceans[i].BottleCount = count
  }

  return OkResponse(c, oceans)
}
