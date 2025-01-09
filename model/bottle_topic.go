package model

// BottleTopic 漂流瓶话题关联表
type BottleTopic struct {
    BaseModel
    TopicID  uint   `gorm:"not null;index;comment:话题ID" json:"topic_id"`
    UserID   uint   `gorm:"not null;index;comment:参与用户ID" json:"user_id"`
    BottleID uint   `gorm:"not null;index;comment:参与漂流瓶ID" json:"bottle_id"`
    // 关联
    Topic  Topic  `gorm:"foreignKey:TopicID" json:"topic"`
    User   User   `gorm:"foreignKey:UserID" json:"user"`
    Bottle Bottle `gorm:"foreignKey:BottleID" json:"bottle"`
}

// TableName 指定表名
func (BottleTopic) TableName() string {
    return "bottle_topics"
}
