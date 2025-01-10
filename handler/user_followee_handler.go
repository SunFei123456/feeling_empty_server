package handler

import (
  "fangkong_xinsheng_app/service"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "github.com/labstack/gommon/log"
  "net/http"
  "strconv"
)

// UserFolloweesHandler 关注
type UserFolloweesHandler struct {
  FolloweesService service.UserFolloweesService
}

// Index 获取当关注列表
func (h UserFolloweesHandler) Index(c echo.Context) error {
  userId := c.Param("id")
  pageParam := c.QueryParam("page")
  page, err := tools.ParsePageAndCheckParam(pageParam)
  if err != nil {
    log.Errorf("Error when parsing page parameter: %v", err)
    return ErrorResponse(c, http.StatusBadRequest, err.Error()+",请检查page参数")
  }
  // 使用服务层获取数据
  data, meta, err := h.FolloweesService.GetFollowees(userId, page)
  if err != nil {
    log.Errorf("Error when getting followees List: %v", err)
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  // c.JSON(http.StatusOK, PagedOkResponse(data, meta))
  return PagedOkResponse(c, data, meta.Count, meta.Page, meta.Items)
}

// FollowUser 关注
func (h UserFolloweesHandler) FollowUser(c echo.Context) error {

  userId := tools.GetUserIDFromContext(c)
  followeeId := c.Param("id")

  // 执行关注操作
  err := h.FolloweesService.FollowUser(strconv.Itoa(int(userId)), followeeId)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  // 返回结果
  return OkResponse(c, "关注成功")
}

// UnfollowUser 取消关注
func (h UserFolloweesHandler) UnfollowUser(c echo.Context) error {

  userId := tools.GetUserIDFromContext(c)
  followeeId := c.Param("id")
  // 执行取消关注操作
  err := h.FolloweesService.UnfollowUser(strconv.Itoa(int(userId)), followeeId)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  // 返回结果
  return OkResponse(c, "取消关注成功")
}

// GetFollowStatus 获取当前用户与指定用户的关注状态
func (h UserFolloweesHandler) GetFollowStatus(c echo.Context) error {
  userId := tools.GetUserIDFromContext(c)
  followeeId := c.Param("id")
  status, err := h.FolloweesService.GetFollowStatus(strconv.Itoa(int(userId)), followeeId)
  if err != nil {
    log.Errorf("Error when getting followee status: %v", err)
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return OkResponse(c, status)
}
