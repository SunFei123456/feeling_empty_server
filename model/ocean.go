package model

// Ocean 海域表
type Ocean struct {
    BaseModel
    Name string `gorm:"size:50;not null;comment:海域名称" json:"name"`
    Body string `gorm:"type:text;not null;comment:描述" json:"body"`
    Bg   string `gorm:"type:text;not null;comment:海域背景" json:"bg"`
    // 虚拟字段
    BottleCount int64 `gorm:"-" json:"bottle_count"`
}

// TableName 指定表名
func (Ocean) TableName() string {
    return "oceans"
} 