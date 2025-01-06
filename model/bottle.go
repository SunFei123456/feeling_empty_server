package model

// Bottle 漂流瓶表
type Bottle struct {
  BaseModel
  UserID     uint   `gorm:"not null;comment:发布用户ID" json:"user_id"`
  Title      string `gorm:"size:50;comment:标题" json:"title"`
  Content    string `gorm:"type:text;comment:内容" json:"content"`
  ImageURL   string `gorm:"size:255;comment:图片URL" json:"image_url"`
  AudioURL   string `gorm:"size:255;comment:音频URL" json:"audio_url"`
  Mood       string `gorm:"size:50;comment:心情" json:"mood"`
  TopicID    *uint  `gorm:"comment:话题ID" json:"topic_id"`
  IsPublic   bool   `gorm:"default:true;comment:是否公开" json:"is_public"`
  Resonances int    `gorm:"default:0;comment:共鸣量" json:"resonances"`
  Favorites  int    `gorm:"default:0;comment:收藏量" json:"favorites"`
  Shares     int    `gorm:"default:0;comment:分享量" json:"shares"`
  Views      int    `gorm:"default:0;comment:浏览量" json:"views"`
  // 关联
  User  User  `gorm:"foreignKey:UserID" json:"user"`
  Topic Topic `gorm:"foreignKey:TopicID" json:"topic,omitempty"`
}

func (Bottle) TableName() string {
  return "bottles"
}
