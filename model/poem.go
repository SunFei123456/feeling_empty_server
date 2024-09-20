package model

// Poem represents the poems table
type Poem struct {
  BaseModel
  Title      string `gorm:"not null" json:"title"`   // Title
  Content    string `gorm:"not null" json:"content"` // Content
  Author     string `gorm:"not null" json:"author"`  // Author
  Tags       string `json:"tags"`                    // Tags
  CoverImage string `json:"coverImage"`              // Cover Image
}
