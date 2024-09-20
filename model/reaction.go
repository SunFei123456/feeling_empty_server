package model

type Reaction struct {
  BaseModel
  ReactionType   string `json:"reaction_type"` // 反应类型，例如 'like', 'dislike' 等
  TargetableType string `json:"likeable_type"` // 反应对象的类型（如 'article', 'comment'）
  TargetableID   int    `json:"targetable_id"` // 反应对象的 ID（如文章 ID 或评论 ID）
  VisitorID      string `json:"visitor_id"`    // 游客的唯一标识
}
