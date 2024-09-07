package model

import "time"

// Progress 表示进度记录模型
type Progress struct {
  BaseModel
  Date                   time.Time `gorm:"not null" json:"date"`                         // 记录日期
  Username               string    `gorm:"size:255;not null" json:"username"`            // 用户名
  Title                  string    `gorm:"size:255;not null" json:"title"`               // 记录标题
  Description            *string   `gorm:"size:255" json:"description"`                  // 记录的简要描述（可选）
  Body                   *string   `gorm:"type:text" json:"body"`                        // 详细记录内容（可选）
  ExpectedCompletionDays *int      `gorm:"default:NULL" json:"expected_completion_days"` // 预计完成天数（可选）
  Times                  int       `gorm:"default:0" json:"times"`                       // 记录次数，默认为 0
}
