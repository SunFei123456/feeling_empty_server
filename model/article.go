package model

type Article struct {
  BaseModel
  UserInfo   User   `gorm:"foreignKey:UserID" json:"user"`
  Title      string `gorm:"size:255;not null" json:"title"`
  Content    string `gorm:"type:text;not null" json:"content"`
  CoverImage string `gorm:"size:255" json:"coverImage"`
  UserID     int    `gorm:"not null" json:"userID"`
  Tags       string `gorm:"type:json" json:"tags"`
  Category   string `gorm:"size:50" json:"category"`
  // 是否为热门帖子
  IsHot bool `gorm:"default:false" json:"isHot"`
  // 浏览量
  Views int `gorm:"default:0" json:"views"`
  // 评论
  Comments []Comment `gorm:"foreignKey:commentable_id" json:"comments"` // 指定外键

  //// 点赞数
  //Likes int `gorm:"default:0" json:"likes"`
  //// 评论数
  //Comments int `gorm:"default:0" json:"comments"`

}
