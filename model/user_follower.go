package model

import (
  "time"
)

type UserFollower struct {
  BaseModel
  Follower     User
  Followee     User
  FollowerId   string    `gorm:"" json:"follower_id"`
  FolloweeId   string    `gorm:"" json:"followee_id"`
  FollowAt     time.Time `json:"follow_at"`
  FollowStatus string    `json:"follow_status" gorm:"->"` // 关注状态
}
