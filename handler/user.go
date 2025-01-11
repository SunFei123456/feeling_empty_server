package handler

import (
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/service"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "github.com/labstack/echo/v4"
  "net/http"
)

type UserHandler struct {
  userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
  return &UserHandler{userService: userService}
}

// HandleRegister 处理用户注册
func (h *UserHandler) HandleRegister(c echo.Context) error {
  var req structs.RegisterRequest
  if err := c.Bind(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "请求参数格式错误"+err.Error())
  }

  if err := c.Validate(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "参数校验错误"+err.Error())
  }

  user := &model.User{
    Email:    req.Email,
    Password: req.Password,
    Phone:    req.Phone,
    Nickname: req.Nickname,
  }

  if err := h.userService.Register(user); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  // 生成JWT token
  token, err := tools.GenerateJWTToken(user.ID)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "生成token失败")
  }

  // 返回token和用户信息
  return OkResponse(c, map[string]interface{}{
    "token": token,
    "user":  user,
  })
}

// HandleLogin 处理用户登录
func (h *UserHandler) HandleLogin(c echo.Context) error {
  var req structs.LoginRequest
  if err := c.Bind(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数"+err.Error())
  }

  if err := c.Validate(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "参数校验错误"+err.Error())
  }

  user, err := h.userService.Login(req.Account, req.Password)
  if err != nil {
    return ErrorResponse(c, http.StatusUnauthorized, "登录失败"+err.Error())
  }

  // 生成JWT token
  token, err := tools.GenerateJWTToken(user.ID)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "生成token失败"+err.Error())
  }

  return OkResponse(c, map[string]interface{}{
    "token": token,
    "user":  user,
  })
}

// HandleGetCurrentUser 获取当前用户信息
func (h *UserHandler) HandleGetCurrentUser(c echo.Context) error {
  userID := tools.GetUserIDFromContext(c)
  user, err := h.userService.GetUserByID(userID)
  if err != nil {
    return ErrorResponse(c, http.StatusNotFound, "User not found")
  }

  return OkResponse(c, tools.ToMap(user, "id", "email", "nickname", "avatar", "sex"))
}

// HandleUpdateCurrentUser 更新当前用户信息
func (h *UserHandler) HandleUpdateCurrentUser(c echo.Context) error {
  var req structs.UpdateUserRequest
  if err := c.Bind(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数")
  }

  if err := c.Validate(&req); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, err.Error())
  }

  userID := tools.GetUserIDFromContext(c)
  updates := make(map[string]interface{})

  // 只更新非空字段
  if req.Nickname != "" {
    updates["nickname"] = req.Nickname
  }
  if req.Avatar != "" {
    updates["avatar"] = req.Avatar
  }
  if req.Sex != nil {
    updates["sex"] = *req.Sex
  }

  if err := h.userService.UpdateUser(userID, updates); err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "更新用户信息失败: "+err.Error())
  }

  // 获取更新后的用户信息
  user, err := h.userService.GetUserByID(userID)
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, "获取更新后的用户信息失败")
  }

  return OkResponse(c, map[string]interface{}{
    "message": "用户信息更新成功",
    "user":    user,
  })
}

// HandleSendEmailCode 处理发送邮箱验证码
func (h *UserHandler) HandleSendEmailCode(c echo.Context) error {
    var req structs.SendEmailCodeRequest
    if err := c.Bind(&req); err != nil {
        return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数"+err.Error())
    }

    if err := c.Validate(&req); err != nil {
        return ErrorResponse(c, http.StatusBadRequest, "参数校验错误"+err.Error())
    }

    if err := h.userService.SendEmailCode(req.Email); err != nil {
        return ErrorResponse(c, http.StatusInternalServerError, err.Error())
    }

    return OkResponse(c, "验证码已发送")
}

// HandleQQEmailLogin 处理QQ邮箱验证码登录
func (h *UserHandler) HandleQQEmailLogin(c echo.Context) error {
    var req structs.QQEmailLoginRequest
    if err := c.Bind(&req); err != nil {
        return ErrorResponse(c, http.StatusBadRequest, "无效的请求参数"+err.Error())
    }

    if err := c.Validate(&req); err != nil {
        return ErrorResponse(c, http.StatusBadRequest, "参数校验错误"+err.Error())
    }

    user, err := h.userService.LoginWithEmailCode(req.Email, req.Code)
    if err != nil {
        return ErrorResponse(c, http.StatusUnauthorized, err.Error())
    }

    // 生成JWT token
    token, err := tools.GenerateJWTToken(user.ID)
    if err != nil {
        return ErrorResponse(c, http.StatusInternalServerError, "生成token失败"+err.Error())
    }

    return OkResponse(c, map[string]interface{}{
        "token": token,
        "user":  tools.ToMap(user, "id", "email", "nickname", "avatar", "sex"),
    })
}
