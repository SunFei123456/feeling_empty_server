package db

import (
  "fangkong_xinsheng_app/model"
  "gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库结构
func AutoMigrate(db *gorm.DB) error {
  // 修改 users 表的 password 字段长度
  err := db.Exec("ALTER TABLE users MODIFY COLUMN password VARCHAR(255) NOT NULL").Error
  if err != nil {
    return err
  }

  // 自动迁移其他表结构
  return db.AutoMigrate(
    &model.User{},
  )
}
