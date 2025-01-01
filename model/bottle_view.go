package model

// BottleView 用户打开的漂流瓶记录表
type BottleView struct {
  BaseModel
  BottleID uint `gorm:"not null;comment:漂流瓶ID" json:"bottle_id"`
  UserID   uint `gorm:"not null;comment:用户ID" json:"user_id"`
  // 关联
  Bottle Bottle `gorm:"foreignKey:BottleID" json:"bottle"`
  User   User   `gorm:"foreignKey:UserID" json:"user"`
}

func (BottleView) TableName() string {
  return "bottle_views"
}
