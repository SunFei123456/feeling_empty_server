package structs

// TopicQueryParams 话题查询参数
type TopicQueryParams struct {
  Page     int `query:"page" validate:"min=1"`
  PageSize int `query:"page_size" validate:"min=1,max=50"`
}

// CreateTopicRequest 创建话题请求
type CreateTopicRequest struct {
  Title   string `json:"title" validate:"required,min=2,max=50"`
  Type    int    `json:"type" validate:"required,oneof=0 1"`
  BgImage string `json:"bg_image" validate:"omitempty,url"`
}
