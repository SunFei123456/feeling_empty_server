package main

import (
  "fmt"
  "github.com/joho/godotenv"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "gorm.io/gorm/logger"
  "log"
  "os"
  "quick-start/db"
  "quick-start/router"
)

func main() {
  // 加载环境变量
  fmt.Println("你好")
  loadEnv()
  // 初始化数据库连接
  initDataBases()
  // 启动 HTTP 服务器
  startServer()
  fmt.Println("哈哈哈")
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
  // 获取.env的DSN变量
  dsn := os.Getenv("DSN")
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

func startServer() {
  // 初始化 Echo 实例
  e := echo.New()
  e.Use(middleware.CORS())
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Static("/static", "static")
  // 引入注册路由
  router.SetupRoutes(e)
  // 启动服务器
  e.Logger.Fatal(e.Start(":8080"))
}
