package service

import (
  "context"
  "errors"
  "fangkong_xinsheng_app/db"
  "fangkong_xinsheng_app/model"
  "fangkong_xinsheng_app/tools"
  "fmt"
  "golang.org/x/crypto/bcrypt"
  "gorm.io/gorm"
  "math/rand"
  "strings"
  "time"
)

type UserService struct {
  db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
  return &UserService{db: db}
}

// Register 用户注册
func (s *UserService) Register(user *model.User) error {
  // 检查邮箱是否已存在
  var count int64
  s.db.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
  if count > 0 {
    return errors.New("邮箱已存在")
  }

  // 如果提供了手机号，检查手机号是否已存在
  if user.Phone != "" {
    s.db.Model(&model.User{}).Where("phone = ?", user.Phone).Count(&count)
    if count > 0 {
      return errors.New("手机号已存在")
    }
  }

  // 设置默认昵称（如果没有提供）
  if user.Nickname == "" {
    user.Nickname = "用户" + user.Email[:5] // 使用邮箱前5位作为默认昵称
  }

  // 密码加密前的验证
  if len(user.Password) < 6 {
    return errors.New("密码长度不能小于6位")
  }

  // 密码加密前先打印日志
  fmt.Printf("注册时的原始密码: %s\n", user.Password)

  // 使用固定的 cost 值
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
  if err != nil {
    return fmt.Errorf("密码加密失败: %v", err)
  }

  // 验证生成的哈希是否正确
  if !strings.HasPrefix(string(hashedPassword), "$2a$") {
    return errors.New("密码加密结果格式错误")
  }

  user.Password = string(hashedPassword)
  fmt.Printf("注册时生成的加密密码: %s\n", user.Password)

  // 创建用户
  if err := s.db.Create(user).Error; err != nil {
    return fmt.Errorf("创建用户失败: %v", err)
  }

  // 验证是否正确保存
  var savedUser model.User
  if err := s.db.Where("email = ?", user.Email).First(&savedUser).Error; err != nil {
    return fmt.Errorf("验证用户创建失败: %v", err)
  }

  if !strings.HasPrefix(savedUser.Password, "$2a$") {
    return errors.New("保存的密码哈希格式错误")
  }

  return nil
}

// Login 用户登录
func (s *UserService) Login(account, password string) (*model.User, error) {
  var user model.User

  // 先尝试用邮箱登录
  err := s.db.Where("email = ?", account).First(&user).Error
  if err != nil {
    // 如果邮箱登录失败，尝试用手机号登录
    err = s.db.Where("phone = ?", account).First(&user).Error
    if err != nil {
      return nil, errors.New("账号不存在")
    }
  }

  // 验证密码
  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
  if err != nil {
    return nil, errors.New("密码错误")
  }
  return &user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
  var user model.User
  if err := s.db.First(&user, id).Error; err != nil {
    return nil, err
  }
  return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
  // 先检查用户是否存在
  var user model.User
  if err := s.db.First(&user, id).Error; err != nil {
    return fmt.Errorf("用户不存在: %v", err)
  }

  // 执行更新操作
  result := s.db.Model(&model.User{}).
    Where("id = ?", id).
    Updates(updates)

  if result.Error != nil {
    return fmt.Errorf("更新失败: %v", result.Error)
  }

  // 检查是否有记录被更新
  if result.RowsAffected == 0 {
    return fmt.Errorf("没有记录被更新")
  }

  return nil
}

// SendEmailCode 发送邮箱验证码
func (s *UserService) SendEmailCode(email string) error {
  // 验证是否是QQ邮箱
  if !strings.HasSuffix(email, "@qq.com") {
    return errors.New("请使用QQ邮箱")
  }

  // 生成6位随机验证码
  code := fmt.Sprintf("%06d", rand.Intn(1000000))

  // 将验证码保存到Redis，设置5分钟过期
  key := fmt.Sprintf("email_code:%s", email)
  if err := db.RedisClient.Set(context.Background(), key, code, 5*time.Minute).Err(); err != nil {
    return fmt.Errorf("保存验证码失败: %v", err)
  }

  // 发送邮件
  subject := "放空心声 - 验证码"
  body := fmt.Sprintf("您的验证码是：%s，5分钟内有效。", code)
  if err := tools.SendEmail(email, subject, body); err != nil {
    return fmt.Errorf("发送邮件失败: %v", err)
  }

  return nil
}

// LoginWithEmailCode QQ邮箱验证码登录
func (s *UserService) LoginWithEmailCode(email, code string) (*model.User, error) {
  ctx := context.Background()
  // 验证是否是QQ邮箱
  if !strings.HasSuffix(email, "@qq.com") {
    return nil, errors.New("请使用QQ邮箱")
  }

  // 从Redis获取验证码
  key := fmt.Sprintf("email_code:%s", email)
  savedCode, err := db.RedisClient.Get(ctx, key).Result()
  if err != nil {
    return nil, errors.New("验证码已过期或不存在")
  }

  // 验证码校验
  if code != savedCode {
    return nil, errors.New("验证码错误")
  }

  // 删除已使用的验证码
  db.RedisClient.Del(ctx, key)

  // 查找用户
  var user model.User
  err = s.db.Where("email = ?", email).First(&user).Error
  if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, fmt.Errorf("用户查询失败: %v", err)
  }

  // 如果用户不存在，则创建新用户
  if errors.Is(err, gorm.ErrRecordNotFound) {
    user = model.User{
      Email:    email,
      Nickname: "用户" + strings.Split(email, "@")[0], // 使用QQ号作为默认昵称
    }

    err = s.db.Create(&user).Error
    if err != nil {
      return nil, fmt.Errorf("用户创建失败: %v", err)
    }
  }

  return &user, nil
}
