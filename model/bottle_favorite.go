package model

// BottleFavorite 用户收藏漂流瓶表
type BottleFavorite struct {
    BaseModel
    UserID   uint   `gorm:"not null;index;comment:用户ID" json:"user_id"`
    BottleID uint   `gorm:"not null;index;comment:漂流瓶ID" json:"bottle_id"`
    // Note     string `gorm:"type:varchar(255);comment:收藏备注" json:"note"`
    // 关联
    User   User   `gorm:"foreignKey:UserID" json:"user"`
    Bottle Bottle `gorm:"foreignKey:BottleID" json:"bottle"`
}

// TableName 指定表名
func (BottleFavorite) TableName() string {
    return "bottle_favorites"
} 