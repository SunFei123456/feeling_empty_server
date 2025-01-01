package model

// User 用户表
type User struct {
  BaseModel
  Nickname string `gorm:"size:50;default:'';comment:昵称" json:"nickname"`
  Avatar   string `gorm:"size:255;default:'';comment:头像" json:"avatar"`
  Sex      int8   `gorm:"default:0;comment:性别 0-未知 1-男 2-女" json:"sex"`
  Email    string `gorm:"size:100;uniqueIndex;not null;comment:邮箱" json:"email"`
  Password string `gorm:"size:255;not null;comment:密码" json:"-"`
  Phone    string `gorm:"size:20;uniqueIndex:idx_phone;default:null;comment:手机号" json:"phone"`
}

// TableName 指定表名
func (User) TableName() string {
  return "users"
}
