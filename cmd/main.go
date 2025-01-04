package main

import (
  "fangkong_xinsheng_app/db"
  "fangkong_xinsheng_app/router"
  "fangkong_xinsheng_app/tools"
  "fmt"
  "github.com/joho/godotenv"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
  "log"
  "os"
)

func main() {
  // 加载环境变量
  loadEnv()
  // 初始化数据库连接
  initDataBases()
  // 执行数据库迁移
  //if err := db.AutoMigrate(db.DB); err != nil {
  //  log.Fatalf("数据库迁移失败: %v", err)
  //}
  // 初始化 Echo 实例
  e := echo.New()
  // 注册验证器
  e.Validator = tools.NewCustomValidator()
  e.Use(middleware.CORS())
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Static("/static", "static")
  // 引入注册路由
  router.SetupRoutes(e)
  // 启动服务器
  print("服务器允许在 http://localhost:8080 访问")
  e.Logger.Fatal(e.Start(":8080"))
}

func loadEnv() {
  env := os.Getenv("ENV")
  if env != "production" { // 生产环境是用docker运行的，会用--env-file参数指定.env文件，不需要手动加载
    err := godotenv.Load()
    if err != nil {
      log.Fatalf("Error loading .env file: %v", err)
    }
  }
}

func initDataBases() {
  // 设置Gorm Logger
  newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
      SlowThreshold: 200, // 慢 SQL 阈值，单位毫秒
      LogLevel:      logger.Info,
      Colorful:      true,
    },
  )
  // 获取.env的DSN变量 如果本地, 则DSN, 远端 则REMOTE_DSN
  dsn := ""
  if os.Getenv("ENV") == "production" {
    dsn = os.Getenv("REMOTE_DSN")
  } else {
    dsn = os.Getenv("DSN")
  }
  // 调用 Open 方法，传入驱动名和连接字符串
  var err error
  db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: newLogger,
  })
  // 检查是否有错误
  if err != nil {
    fmt.Println("连接数据库失败：", err)
    return
  }

  // 打印成功信息
  fmt.Println("连接数据库成功", db.DB)
}
