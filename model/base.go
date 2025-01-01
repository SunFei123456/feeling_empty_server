package model

import (
  "time"
)

// BaseModel 基础模型
type BaseModel struct {
  ID        uint      `gorm:"primarykey;autoIncrement;type:bigint unsigned" json:"id"`
  CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
  UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
