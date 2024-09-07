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
}
