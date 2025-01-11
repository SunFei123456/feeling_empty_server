package structs

// RegisterRequest 注册请求
type RegisterRequest struct {
  Email    string `json:"email" validate:"required,email"`
  Password string `json:"password" validate:"required,min=6,max=20"`
  Phone    string `json:"phone" validate:"omitempty,len=11"`
  Nickname string `json:"nickname" validate:"omitempty,min=2,max=50"`
}

// LoginRequest 登录请求
type LoginRequest struct {
  // 账号可以是邮箱或手机号
  Account  string `json:"account" validate:"required"`
  Password string `json:"password" validate:"required"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
  Nickname string `json:"nickname" validate:"omitempty,min=2,max=50"`
  Avatar   string `json:"avatar" validate:"omitempty,url"`
  Sex      *int8  `json:"sex" validate:"omitempty,oneof=0 1 2"`
}

// UserResponse 用户信息响应
type UserResponse struct {
  ID       uint   `json:"id"`
  Nickname string `json:"nickname"`
  Avatar   string `json:"avatar"`
  Sex      int8   `json:"sex"`
  Email    string `json:"email"`
  Phone    string `json:"phone"`
}

// SendEmailCodeRequest 发送邮箱验证码请求
type SendEmailCodeRequest struct {
    Email string `json:"email" validate:"required,email"`
}

// QQEmailLoginRequest QQ邮箱验证码登录请求
type QQEmailLoginRequest struct {
    Email string `json:"email" validate:"required,email"`
    Code  string `json:"code" validate:"required,len=6"`
}
