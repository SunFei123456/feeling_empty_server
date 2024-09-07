package model

import "time"

// 定义一个BaseModel结构体
type BaseModel struct {
  ID        int       `gorm:"primary_key" json:"id"`
  CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
  UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
