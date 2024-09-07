package handle

import (
  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
  "net/http"
  "quick-start/db"
  "quick-start/model"
)

type ProgressHandler struct{}

// 获取进度记录
func (h *ProgressHandler) Get(c echo.Context) error {
  var progress []model.Progress
  // 使用 Find 查询所有记录
  if err := db.DB.Table("progress_records").Order("created_at desc").Find(&progress).Error; err != nil {
    // 如果错误不是由于没有记录，而是其他数据库错误，则返回 500 错误
    if err == gorm.ErrRecordNotFound {
      return c.JSON(http.StatusNotFound, echo.Map{"error": "数据不存在"})
    }
    return c.JSON(http.StatusInternalServerError, echo.Map{"error": "查询失败", "details": err.Error()})
  }
  // 返回查询到的数据
  return c.JSON(http.StatusOK, echo.Map{"data": progress})
}
