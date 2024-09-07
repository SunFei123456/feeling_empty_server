package model

// 定义一个User 结构体,用来表示user表
// User 结构体表示用户表
type User struct {
  BaseModel
  Nickname string `json:"nickname"`
  Gender   string `json:"gender"`
  Avatar   string `json:"avatar"`
  Email    string `json:"email"`
  Password string `json:"password"`
  Phone    string `json:"phone"`
  Address  string `json:"address"`
}
