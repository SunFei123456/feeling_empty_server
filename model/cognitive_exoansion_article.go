package model

type CognitiveExpansionArticle struct {
  BaseModel
  Title      string `gorm:"size:255;not null" json:"title"`
  Content    string `gorm:"type:text;not null" json:"body"`
  CoverImage string `gorm:"size:255" json:"cover_image"`
  Tags       string `gorm:"size:255" json:"tags"`
}
