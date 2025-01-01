package structs

import (
  "errors"
)

// CreateBottleRequest 创建漂流瓶请求
type CreateBottleRequest struct {
  Content  string `json:"content" validate:"omitempty,max=1000"`
  ImageURL string `json:"image_url" validate:"omitempty,url"`
  AudioURL string `json:"audio_url" validate:"omitempty,url"`
  Mood     string `json:"mood" validate:"required"`
  TopicID  *uint  `json:"topic_id"`
  IsPublic bool   `json:"is_public" validate:"required"`
}

// Validate 自定义验证函数
func (r *CreateBottleRequest) Validate() error {
  // 至少需要有一个内容（文字、图片或音频）
  if r.Content == "" && r.ImageURL == "" && r.AudioURL == "" {
    return errors.New("漂流瓶必须包含文字、图片或音频内容")
  }

  // 图片和音频不能同时存在
  if r.ImageURL != "" && r.AudioURL != "" {
    return errors.New("图片和音频不能同时存在")
  }

  return nil
}

// UpdateBottleRequest 更新漂流瓶请求
type UpdateBottleRequest struct {
  Content  string `json:"content" validate:"omitempty,max=1000"`
  ImageURL string `json:"image_url" validate:"omitempty,url"`
  AudioURL string `json:"audio_url" validate:"omitempty,url"`
  Mood     string `json:"mood"`
  IsPublic *bool  `json:"is_public"`
}

// BottleQueryParams 漂流瓶查询参数
type BottleQueryParams struct {
  UserID   uint   `query:"user_id"`
  TopicID  uint   `query:"topic_id"`
  IsPublic *bool  `query:"is_public"`
  Page     int    `query:"page" validate:"min=1"`
  PageSize int    `query:"page_size" validate:"min=1,max=50"`
  Sort     string `query:"sort"`
}

// HotBottleQueryParams 热门漂流瓶查询参数
type HotBottleQueryParams struct {
  TimeRange string `query:"time_range" validate:"omitempty,oneof=day week month all"` // 时间范围：day-24小时内 week-一周内 month-一个月内 all-全部
  Page      int    `query:"page" validate:"min=1"`
  PageSize  int    `query:"page_size" validate:"min=1,max=50"`
}
