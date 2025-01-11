package service

import (
  "fangkong_xinsheng_app/db"
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/structs"
  "fangkong_xinsheng_app/tools"
  "fmt"
  "github.com/labstack/gommon/log"
  "gorm.io/gorm"
  "time"
)

// 关注
type UserFolloweesService struct {
}

// GetFollowees 获取关注的人员列表
func (s UserFolloweesService) GetFollowees(userId string, page int) ([]map[string]any, structs.Pagination, error) {
  // 默认每页显示 20 条
  itemsPerPage := 20
  // 计算当前页的起始位置
  offset := (page - 1) * itemsPerPage

  // 获取关注列表人的总数
  var totalCount int64
  err := db.DB.Table("user_followers").Where("followee_id = ?", userId).Count(&totalCount).Error
  if err != nil {
    log.Errorf("GetFollowees error: %v", err)
    return nil, structs.Pagination{}, err
  }

  var results []model.UserFollower

  if err := db.DB.
      Select(`
    user_followers.followee_id,
    user_followers.follow_at,
          CASE
            WHEN user_followers.follower_id = ? AND uf2.follower_id IS NOT NULL THEN 'mutual_following'  -- 互相关注
            WHEN user_followers.follower_id = ? THEN 'following'                                        -- 我关注了他
            WHEN uf2.follower_id = ? THEN 'followed'                                         -- 他关注了我
            ELSE 'not_following'                                                             -- 未关注
          END AS follow_status
        `, userId, userId, userId).
    Joins("JOIN users ON users.id = user_followers.followee_id").
    Joins("LEFT JOIN user_followers AS uf2 ON uf2.follower_id = user_followers.followee_id AND uf2.followee_id = ?", userId).
    Preload("Followee", func(db *gorm.DB) *gorm.DB { return db.Select("id", "nickname", "avatar", "sex") }).
    Where("user_followers.follower_id = ?", userId).
    Order("follow_at DESC").
    Limit(itemsPerPage).
    Offset(offset).
    Find(&results).Error; err != nil {
    log.Errorf("GetFollowees error: %v", err)
    return nil, structs.Pagination{}, err
  }

  // 准备一个切片,用于存储返回的数据 [{},{},{},{}...]
  responses := make([]map[string]any, len(results))
  for i, result := range results {
    responses[i] = map[string]any{
      "followee_id":   result.FolloweeId,
      "follow_at":     result.FollowAt,
      "user":          tools.ToMap(result.Followee, "id", "nickname", "avatar", "sex"),
      "follow_status": result.FollowStatus,
    }
  }
  // 计算分页信息
  meta := structs.NewPagination(totalCount, page, itemsPerPage)
  // 返回结果
  return responses, meta, nil
}

// FollowUser 关注
func (s UserFolloweesService) FollowUser(followerId, followeeId string) error {
  followerData := &model.UserFollower{
    FollowerId: followerId,
    FolloweeId: followeeId,
    FollowAt:   time.Now(),
  }
  err := db.DB.Table("user_followers").Create(followerData).Error
  if err != nil {
    log.Errorf("关注失败: %v", err)
    return fmt.Errorf("关注失败: %v", err)
  }
  return nil
}

// UnfollowUser 取消关注
func (s UserFolloweesService) UnfollowUser(followerId, followeeId string) error {
  // 已经关注，执行取关操作
  err := db.DB.Where("follower_id = ? AND followee_id = ?", followerId, followeeId).
    Delete(&model.UserFollower{}).Error
  if err != nil {
    log.Errorf("取消关注失败: %v", err)
    return fmt.Errorf("取消关注失败: %v", err)
  }
  return err
}

// GetFollowStatus 获取当前用户与指定用户之间的关注状态
func (s UserFolloweesService) GetFollowStatus(followerId, followeeId string) (string, error) {
  var followerCount, followeeCount int64

  // 事务
  err := db.DB.Transaction(func(tx *gorm.DB) error {
    // 查询当前用户对指定用户的关注记录
    if err := tx.Model(&model.UserFollower{}).
      Where("follower_id = ? AND followee_id = ?", followerId, followeeId).
      Count(&followerCount).Error; err != nil {
      return err
    }

    // 查询指定用户对当前用户的关注记录
    if err := tx.Model(&model.UserFollower{}).
      Where("follower_id = ? AND followee_id = ?", followeeId, followerId).
      Count(&followeeCount).Error; err != nil {
      return err
    }

    return nil
  })

  if err != nil {
    return "", err
  }

  // 根据查询结果返回对应的状态
  switch {
  case followerCount == 0 && followeeCount == 0: // 我没关注他,他没关注我
    return "not_following", nil
  case followerCount > 0 && followeeCount == 0: // 我关注了他，但他没有关注我
    return "following", nil
  case followerCount == 0 && followeeCount > 0: // 我没关注他，他关注了我
    return "followed", nil
  case followerCount > 0 && followeeCount > 0: // 互相关注
    return "mutual_following", nil
  default:
    return "", gorm.ErrRecordNotFound
  }
}
