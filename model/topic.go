package model

type Topic struct {
  BaseModel
  Title   string `gorm:"not null;comment:标题" json:"title"`
  Desc    string `gorm:"type:text;comment:描述" json:"desc"`
  Type    int    `gorm:"not null;comment:类型" json:"type"` // 0: 用户, 1: 系统
  Views   int    `gorm:"default:0;comment:浏览量" json:"views"`
  BgImage string `gorm:"type:text;comment:背景" json:"bg_image"`
}

// TableName 指定表名
func (Topic) TableName() string {
  return "topics"
}
