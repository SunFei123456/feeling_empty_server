package model

// BottleResonance 用户共振漂流瓶表
type BottleResonance struct {
    BaseModel
    UserID      uint   `gorm:"not null;index;comment:用户ID" json:"user_id"`
    BottleID    uint   `gorm:"not null;index;comment:漂流瓶ID" json:"bottle_id"`
    // 关联
    User          User   `gorm:"foreignKey:UserID" json:"user"`
    Bottle        Bottle `gorm:"foreignKey:BottleID" json:"bottle"`
}

// TableName 指定表名
func (BottleResonance) TableName() string {
    return "bottle_resonances"
} 