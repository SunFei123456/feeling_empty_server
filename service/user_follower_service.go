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

// 粉丝
type UserFollowersService struct{}

// GetFollowersWithFilter 根据过滤条件获取粉丝列表
func (s UserFollowersService) GetFollowersWithFilter(userId string, followedAfter time.Time, page, itemsPerPage int) ([]map[string]any, structs.Pagination, error) {
  // 计算当前页的起始位置
  offset := (page - 1) * itemsPerPage

  // 计算粉丝总数
  var totalCount int64

  // 如果followedAfter 不为空
  if !time.Time.IsZero(followedAfter) {
    err := db.DB.Model(&model.UserFollower{}).Where("followee_id = ? AND follow_at >= ?", userId, followedAfter).Count(&totalCount).Error
    if err != nil {
      log.Errorf("Error when counting followers: %v", err)
      return nil, structs.Pagination{}, err
    }
  }
  err := db.DB.Model(&model.UserFollower{}).Where("followee_id = ?", userId).Count(&totalCount).Error
  if err != nil {
    log.Errorf("Error when counting followers: %v", err)
    return nil, structs.Pagination{}, err
  }

  var results []model.UserFollower
  // 构建查询
  query := db.DB.
    Preload("Follower", func(db *gorm.DB) *gorm.DB { return db.Select("id", "nickname", "avatar", "sex") }).
      Select(`
            user_followers.follower_id,
            user_followers.follow_at,
            CASE
                WHEN user_followers.followee_id = ? AND uf2.followee_id IS NOT NULL THEN 'mutual_following'  -- 互相关注
                WHEN user_followers.followee_id = ? THEN 'followed'                                         -- 他关注了我
                WHEN uf2.followee_id = ? THEN 'following'                                                   -- 我关注了他
                ELSE 'not_following'                                                                        -- 未关注
            END AS follow_status
        `, userId, userId, userId).
    Joins("JOIN users ON users.id = user_followers.follower_id").
    Joins("LEFT JOIN user_followers AS uf2 ON uf2.followee_id = user_followers.follower_id AND uf2.follower_id = ?", userId).
    Where("user_followers.followee_id = ?", userId).
    Order("follow_at DESC").
    Limit(itemsPerPage).
    Offset(offset)
  // 如果有过滤条件，则应用过滤条件, 时间 为 空 true 就说明 无过滤条件,不为空 说明有过滤条件
  if !time.Time.IsZero(followedAfter) {
    // 进行解析，只取到日期部分
    date := followedAfter.Format("2006-01-02") // 格式化为日期字符串
    query = query.Where("user_followers.follow_at >= ?", date)
  }

  // 执行查询
  if err := query.Find(&results).Error; err != nil {
    log.Errorf("Error when getting followers List: %v", err)
    return nil, structs.Pagination{}, fmt.Errorf("error when getting followers List: %v", err)
  }

  // 构建返回结果
  responses := make([]map[string]any, len(results))
  for i, result := range results {
    responses[i] = tools.ToMap(result, "follow_at", "follow_status", "follower_id")
    responses[i]["user"] = tools.ToMap(result.Follower, "id", "nickname", "avatar", "sex")
  }

  meta := structs.NewPagination(totalCount, page, itemsPerPage)
  return responses, meta, nil
}

// GetFollowers 获取粉丝列表
func (s UserFollowersService) GetFollowers(userId string, page, itemsPerPage int) ([]map[string]any, structs.Pagination, error) {
  return s.GetFollowersWithFilter(userId, time.Time{}, page, itemsPerPage)
}

// GetRecentThreeDaysFansList 获取过去三天内新增的粉丝列表
func (s UserFollowersService) GetRecentThreeDaysFansList(userId string, page, itemsPerPage int) ([]map[string]any, structs.Pagination, error) {
  // 计算三天前的时间
  threeDaysAgo := time.Now().AddDate(0, 0, -3)
  return s.GetFollowersWithFilter(userId, threeDaysAgo, page, itemsPerPage)
}
