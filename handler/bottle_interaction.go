package handler

import (
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "fmt"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "strconv"
)

type BottleInteractionHandler struct {
  db *gorm.DB
}

func NewBottleInteractionHandler(db *gorm.DB) *BottleInteractionHandler {
  return &BottleInteractionHandler{db: db}
}

// HandleResonateBottle 共振漂流瓶
func (h *BottleInteractionHandler) HandleResonateBottle(c echo.Context) error {
  userID := tools.GetUserIDFromContext(c)

  // 1. 获取并验证 bottle_id
  bottleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的漂流瓶ID")
  }

  // 2. 检查漂流瓶是否存在
  var bottle model.Bottle
  if err := h.db.First(&bottle, bottleID).Error; err != nil {
    return ErrorResponse(c, http.StatusNotFound, "漂流瓶不存在")
  }

  // 3. 使用事务确保数据一致性
  err = h.db.Transaction(func(tx *gorm.DB) error {
    // 检查是否已经共振过
    var exists bool
    if err := tx.Model(&model.BottleResonance{}).
      Where("user_id = ? AND bottle_id = ?", userID, bottleID).
      Select("count(*) > 0").
      Find(&exists).Error; err != nil {
      return err
    }

    if exists {
      return fmt.Errorf("已经共振过该漂流瓶")
    }

    // 创建共振记录
    resonance := &model.BottleResonance{
      UserID:   userID,
      BottleID: uint(bottleID),
    }
    if err := tx.Create(resonance).Error; err != nil {
      return err
    }

    // 更新漂流瓶共振数
    if err := tx.Model(&bottle).
      UpdateColumn("resonances", gorm.Expr("resonances + ?", 1)).
      Error; err != nil {
      return err
    }

    return nil
  })

  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  return OkResponse(c, "共振成功")
}

// HandleCancelResonateBottle 取消共振漂流瓶
func (h *BottleInteractionHandler) HandleCancelResonateBottle(c echo.Context) error {
  bottleID := c.Param("id")
  userID := tools.GetUserIDFromContext(c)

  err := h.db.Transaction(func(tx *gorm.DB) error {
    // 删除共振记录
    result := tx.Where("user_id = ? AND bottle_id = ?", userID, bottleID).
      Delete(&model.BottleResonance{})

    if result.RowsAffected == 0 {
      return ErrorResponse(c, http.StatusBadRequest, "未共振过该漂流瓶")
    }

    // 更新漂流瓶共振数
    if err := tx.Model(&model.Bottle{}).
      Where("id = ?", bottleID).
      UpdateColumn("resonances", gorm.Expr("resonances - ?", 1)).
      Error; err != nil {
      return ErrorResponse(c, http.StatusInternalServerError, "更新漂流瓶共振数失败: "+err.Error())
    }

    return nil
  })

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "取消共振失败: "+err.Error())
  }

  return OkResponse(c, "取消共振成功")
}

// HandleGetUserResonatedBottles 获取用户共振的漂流瓶列表
func (h *BottleInteractionHandler) HandleGetUserResonatedBottles(c echo.Context) error {
  var params structs.BottleQueryParams
  if err := c.Bind(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
  }

  userID := tools.GetUserIDFromContext(c)

  query := h.db.Model(&model.BottleResonance{}).
    Joins("LEFT JOIN bottles ON bottle_resonances.bottle_id = bottles.id").
    Where("bottle_resonances.user_id = ?", userID).
    Preload("Bottle").
    Preload("Bottle.User")

  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取共振列表失败")
  }

  var resonances []model.BottleResonance
  if err := query.Offset((params.Page - 1) * params.PageSize).
    Limit(params.PageSize).
    Order("bottle_resonances.created_at DESC").
    Find(&resonances).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取共振列表失败")
  }

  var result []map[string]any
  for _, r := range resonances {
    bottleMap := tools.ToMap(&r.Bottle, "id", "title", "content", "image_url", "audio_url", "mood", "topic_id", "created_at", "views", "resonances", "favorites")
    if r.Bottle.User.ID != 0 {
      bottleMap["user"] = tools.ToMap(&r.Bottle.User, "id", "nickname", "avatar", "sex")
    }
    result = append(result, bottleMap)
  }

  return PagedOkResponse(c, result, total, params.Page, params.PageSize)
}

// HandleFavoriteBottle 收藏漂流瓶
func (h *BottleInteractionHandler) HandleFavoriteBottle(c echo.Context) error {
  bottleID := c.Param("id")
  userID := tools.GetUserIDFromContext(c)
  bid, err := strconv.ParseUint(bottleID, 10, 32)
  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的漂流瓶ID")
  }
  err = h.db.Transaction(func(tx *gorm.DB) error {
    // 检查是否已经收藏过
    var exists bool
    if err := tx.Model(&model.BottleFavorite{}).
      Where("user_id = ? AND bottle_id = ?", userID, bottleID).
      Select("count(*) > 0").
      Find(&exists).Error; err != nil {
      return err
    }

    if exists {
      return echo.NewHTTPError(http.StatusBadRequest, "已经收藏过该漂流瓶")
    }

    // 创建收藏记录
    favorite := &model.BottleFavorite{
      UserID:   userID,
      BottleID: uint(bid),
    }
    if err := tx.Create(favorite).Error; err != nil {
      return err
    }

    // 更新漂流瓶收藏数
    if err := tx.Model(&model.Bottle{}).
      Where("id = ?", bottleID).
      UpdateColumn("favorites", gorm.Expr("favorites + ?", 1)).
      Error; err != nil {
      return err
    }

    return nil
  })

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "收藏失败: "+err.Error())
  }

  return OkResponse(c, "收藏成功")
}

// HandleCancelFavoriteBottle 取消收藏漂流瓶
func (h *BottleInteractionHandler) HandleCancelFavoriteBottle(c echo.Context) error {
  bottleID := c.Param("id")
  userID := tools.GetUserIDFromContext(c)

  err := h.db.Transaction(func(tx *gorm.DB) error {
    // 删除收藏记录
    result := tx.Where("user_id = ? AND bottle_id = ?", userID, bottleID).
      Delete(&model.BottleFavorite{})

    if result.RowsAffected == 0 {
      return echo.NewHTTPError(http.StatusBadRequest, "未收藏过该漂流瓶")
    }

    // 更新漂流瓶收藏数
    if err := tx.Model(&model.Bottle{}).
      Where("id = ?", bottleID).
      UpdateColumn("favorites", gorm.Expr("favorites - ?", 1)).
      Error; err != nil {
      return err
    }

    return nil
  })

  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "取消收藏失败: "+err.Error())
  }

  return OkResponse(c, "取消收藏成功")
}

// HandleGetUserFavoriteBottles 获取用户收藏的漂流瓶列表
func (h *BottleInteractionHandler) HandleGetUserFavoriteBottles(c echo.Context) error {
  var params structs.BottleQueryParams
  if err := c.Bind(&params); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
  }

  userID := tools.GetUserIDFromContext(c)

  query := h.db.Model(&model.BottleFavorite{}).
    Joins("LEFT JOIN bottles ON bottle_favorites.bottle_id = bottles.id").
    Where("bottle_favorites.user_id = ?", userID).
    Preload("Bottle").
    Preload("Bottle.User")

  var total int64
  if err := query.Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取收藏列表失败")
  }

  var favorites []model.BottleFavorite
  if err := query.Offset((params.Page - 1) * params.PageSize).
    Limit(params.PageSize).
    Order("bottle_favorites.created_at DESC").
    Find(&favorites).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取收藏列表失败")
  }

  var result []map[string]any
  for _, f := range favorites {
    bottleMap := tools.ToMap(&f.Bottle, "id", "title", "content", "image_url", "audio_url", "mood", "topic_id", "created_at", "views", "resonances", "favorites")
    if f.Bottle.User.ID != 0 {
      bottleMap["user"] = tools.ToMap(&f.Bottle.User, "id", "nickname", "avatar", "sex")
    }
    result = append(result, bottleMap)
  }

  return PagedOkResponse(c, result, total, params.Page, params.PageSize)
}

// HandleShareBottle 分享漂流瓶
func (h *BottleInteractionHandler) HandleShareBottle(c echo.Context) error {
  bottleID := c.Param("id")
  var bottle model.Bottle
  // 更新分享数
  if err := h.db.Model(&bottle).
    Where("id = ?", bottleID).
    UpdateColumn("shares", gorm.Expr("shares + ?", 1)).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "更新分享数失败: "+err.Error())
  }
  // 获取 shares
  if err := h.db.Model(&bottle).
    Where("id = ?", bottleID).
    Select("shares").
    Find(&bottle).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取分享数失败: "+err.Error())
  }

  return OkResponse(c, map[string]any{
    "message": "分享成功",
    "shares":  bottle.Shares,
  })
}

// HandleGetBottleInteractionStatus 获取瓶子的交互状态
func (h *BottleInteractionHandler) HandleGetBottleInteractionStatus(c echo.Context) error {
  bottleID := c.Param("id")
  userID := tools.GetUserIDFromContext(c)

  // 检查当前用户是否收藏了这个瓶子
  var isFavorited bool
  h.db.Model(&model.BottleFavorite{}).
    Where("user_id = ? AND bottle_id = ?", userID, bottleID).
    Select("count(*) > 0").
    Find(&isFavorited)

  // 检查当前用户是否共振了这个瓶子
  var isResonated bool
  h.db.Model(&model.BottleResonance{}).
    Where("user_id = ? AND bottle_id = ?", userID, bottleID).
    Select("count(*) > 0").
    Find(&isResonated)

  return OkResponse(c, map[string]interface{}{
    "is_favorited": isFavorited,
    "is_resonated": isResonated,
  })
}
