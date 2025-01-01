package model

type Topic struct {
  BaseModel
  Title        string `gorm:"not null;comment:标题" json:"title"`
  Description  string `gorm:"type:text;comment:描述" json:"description"`
  Content      string `gorm:"type:text;comment:内容" json:"content"`
  Views        int    `gorm:"default:0;comment:浏览量" json:"views"`
  Participants int    `gorm:"default:0;comment:参与人数" json:"participants"`
  Bottles      int    `gorm:"default:0;comment:漂流瓶数" json:"bottlesint"`
  Status       int    `gorm:"default:0;comment:状态 0-未开始 1-进行中 2-已结束" json:"status"`
  StartDate    string `gorm:"comment:开始时间" json:"start_date"`
  EndDate      string `gorm:"comment:结束时间" json:"end_date"`
  CreatorID    string `gorm:"not null;comment:创建者ID" json:"creator_id"`
  Creator      User   `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}
