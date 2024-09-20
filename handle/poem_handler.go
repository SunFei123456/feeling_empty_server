package handle

import (
  "github.com/labstack/echo/v4"
  "net/http"
  "quick-start/db"
  "quick-start/model"
)

type PoemHandler struct{}

// 获取所有古诗词
func (h *PoemHandler) Get(c echo.Context) error {
  var poems []model.Poem
  err := db.DB.Find(&poems)
  if err.Error != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error.Error())
  }
  return SuccessResponse(c, poems)
}
