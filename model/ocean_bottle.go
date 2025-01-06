package model

// OceanBottle 海域漂流瓶关联表
type OceanBottle struct {
    BaseModel
    OceanID  uint   `gorm:"not null;comment:海域ID" json:"ocean_id"`
    BottleID uint   `gorm:"not null;comment:漂流瓶ID" json:"bottle_id"`
    Ocean    Ocean  `gorm:"foreignKey:OceanID" json:"ocean"`
    Bottle   Bottle `gorm:"foreignKey:BottleID" json:"bottle"`
}

// TableName 指定表名
func (OceanBottle) TableName() string {
    return "ocean_bottles"
} 