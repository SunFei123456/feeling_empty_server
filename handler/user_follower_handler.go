package handler

import (
  "fangkong_xinsheng_app/service"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "github.com/labstack/gommon/log"
  "net/http"
)

// UserFollowersHandler 粉丝
type UserFollowersHandler struct {
  FollowersService service.UserFollowersService
}

// Index getCurrentUserFollowerList 获取粉丝列表
func (h UserFollowersHandler) Index(c echo.Context) error {
  userId := c.Param("id")
  pageParam := c.QueryParam("page")
  page, err := tools.ParsePageAndCheckParam(pageParam)
  if err != nil {
    log.Errorf("Error when parsing page parameter: %v", err)
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }
  itemsPerPage := 20
  // 使用服务层获取数据
  responses, meta, err := h.FollowersService.GetFollowers(userId, page, itemsPerPage)
  if err != nil {
    log.Errorf("Error when getting followers List: %v", err)
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  return PagedOkResponse(c, responses, meta.Count, meta.Page, meta.Items)
}

// GetRecentThreeDaysFansList 获取近三天内新增的粉丝 的列表
func (h UserFollowersHandler) GetRecentThreeDaysFansList(c echo.Context) error {
  userId := c.Param("id")
  pageParam := c.QueryParam("page")
  page, err := tools.ParsePageAndCheckParam(pageParam)
  if err != nil {
    log.Errorf("Error when parsing page parameter: %v", err)
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }
  itemsPerPage := 20
  // 获取数据
  responses, meta, err := h.FollowersService.GetRecentThreeDaysFansList(userId, page, itemsPerPage)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  // 如果responses为空
  if len(responses) == 0 {
    return OkResponse(c, "近三天内没有新增粉丝")
  }
  // 返回数据
  return PagedOkResponse(c, responses, meta.Count, meta.Page, meta.Items)
}
