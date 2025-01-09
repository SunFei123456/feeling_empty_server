package service

import (
    "fangkong_xinsheng_app/model"
    "gorm.io/gorm"
)

type BottleInteractionService struct {
    db *gorm.DB
}

func NewBottleInteractionService(db *gorm.DB) *BottleInteractionService {
    return &BottleInteractionService{db: db}
}

// GetBottleInteractionStatus 获取瓶子的交互状态
func (s *BottleInteractionService) GetBottleInteractionStatus(userID uint, bottleID uint) (bool, bool) {
    var isFavorited bool
    s.db.Model(&model.BottleFavorite{}).
        Where("user_id = ? AND bottle_id = ?", userID, bottleID).
        Select("count(*) > 0").
        Find(&isFavorited)

    var isResonated bool
    s.db.Model(&model.BottleResonance{}).
        Where("user_id = ? AND bottle_id = ?", userID, bottleID).
        Select("count(*) > 0").
        Find(&isResonated)

    return isFavorited, isResonated
}

// EnrichBottleWithInteractionStatus 为瓶子数据添加交互状态
func (s *BottleInteractionService) EnrichBottleWithInteractionStatus(bottleMap map[string]interface{}, userID uint, bottleID uint) {
    isFavorited, isResonated := s.GetBottleInteractionStatus(userID, bottleID)
    bottleMap["is_favorited"] = isFavorited
    bottleMap["is_resonated"] = isResonated
} 