package handle

import (
  "crypto/sha256"
  "errors"
  "fmt"
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "quick-start/db"
  "quick-start/model"
  "strconv"
)

type LikeHandler struct{}

// 点赞或取消点赞文章或评论
func (h *LikeHandler) ToggleLike(c echo.Context) error {
  targetableType := c.Param("type")
  id := c.Param("id")
  reactionType := c.Param("reaction_type")
  idInt, err := strconv.Atoi(id)
  if err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
  }
  visitorID := getVisitorID(c)
  var likeRecord model.Reaction
  if err := db.DB.Where("visitor_id =? AND reaction_type =? AND id =? AND targetable_type =?", visitorID, reactionType, id, targetableType).First(&likeRecord).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      // 检查是否存在相反的反应类型记录，如果有则删除
      if err := db.DB.Where("visitor_id =? AND id =? AND reaction_type!=?", visitorID, id, reactionType).Delete(&model.Reaction{}).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "处理点赞失败"})
      }
      likeRecord = model.Reaction{
        ReactionType: reactionType,
        TargetableID: idInt,
        VisitorID:    visitorID,
      }
      if err := db.DB.Create(&likeRecord).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "点赞失败"})
      }
      return c.JSON(http.StatusOK, echo.Map{"message": "点赞成功"})
    }
  }
  // 有记录，删除
  if err := db.DB.Delete(&likeRecord).Error; err != nil {
    return c.JSON(http.StatusInternalServerError, echo.Map{"error": "取消点赞失败"})
  }
  return c.JSON(http.StatusOK, echo.Map{"message": "取消点赞成功"})
}

// 取消点赞文章或评论
func (h *LikeHandler) Toggledislike(c echo.Context) error {
  reactionType := c.Param("reaction_type")
  id := c.Param("id")
  idInt, err := strconv.Atoi(id)
  if err != nil {
    return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
  }
  visitorID := getVisitorID(c)
  var dislikeRecord model.Reaction
  if err := db.DB.Where("visitor_id =? AND reaction_type =? AND reaction_id =?", visitorID, reactionType, idInt).First(&dislikeRecord).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return c.JSON(http.StatusOK, echo.Map{"message": "没有取消点赞记录"})
    }
  }
  // 有记录，删除
  if err := db.DB.Delete(&dislikeRecord).Error; err != nil {
    return c.JSON(http.StatusInternalServerError, echo.Map{"error": "取消点赞失败"})
  }
  return c.JSON(http.StatusOK, echo.Map{"message": "取消点赞成功"})
}

// 获取游客唯一标识 (可以通过 IP 地址、User-Agent 或者其他方式获取)
func getVisitorID(c echo.Context) string {
  ip := c.RealIP()
  userAgent := c.Request().Header.Get("User-Agent")
  if ip == "" || userAgent == "" {
    return ""
  }
  return generateHash(ip + userAgent)
}

// 生成唯一标识的简单 Hash 函数
func generateHash(input string) string {
  h := sha256.New()
  h.Write([]byte(input))
  return fmt.Sprintf("%x", h.Sum(nil))
}
