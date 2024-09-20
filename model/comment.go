package model

// Comment 结构体映射到数据库中的 Comments 表
type Comment struct {
  BaseModel
  Nickname        string `gorm:"size:100;not null" json:"nickname"`
  Avatar          string `gorm:"size:255" json:"avatar"`
  Body            string `gorm:"type:text;not null" json:"body"`
  Status          string `gorm:"type:enum('visible', 'deleted');default:'visible'" json:"status"`
  VisitorId       string `json:"visitor_id"`
  City            string `json:"city"`
  UserAgent       string `json:"user_agent"`
  CommentableType string `json:"commentable_type"`
  CommentableId   int    `json:"commentable_id"`
}
