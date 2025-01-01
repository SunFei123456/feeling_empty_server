package router

import (
  "fangkong_xinsheng_app/db"
  "fangkong_xinsheng_app/handler"
  "fangkong_xinsheng_app/middleware"
  "fangkong_xinsheng_app/service"
  "github.com/labstack/echo/v4"
)

// SetupRoutes 配置所有路由
func SetupRoutes(e *echo.Echo) {
  // 初始化处理器
  userHandler := handler.NewUserHandler(service.NewUserService(db.DB))
  bottleHandler := handler.NewBottleHandler(db.DB)

  // API 路由组
  api := e.Group("/api/v1")

  // 公开路由组
  auth := api.Group("/auth")
  {
    auth.POST("/register", userHandler.HandleRegister)
    auth.POST("/login", userHandler.HandleLogin)
  }

  // 需要认证的路由组
  authenticated := api.Group("")
  authenticated.Use(middleware.JWT())

  // 用户相关路由
  users := authenticated.Group("/users")
  {
    users.GET("/me", userHandler.HandleGetCurrentUser)
    users.PUT("/me", userHandler.HandleUpdateCurrentUser)
    // 可以添加更多用户相关路由...
  }

  // 漂流瓶相关路由
  bottles := authenticated.Group("/bottles")
  {
    // 基础操作
    bottles.POST("", bottleHandler.HandleCreateBottle)
    bottles.GET("", bottleHandler.HandleGetBottles)
    bottles.GET("/:id", bottleHandler.HandleGetBottle)
    bottles.PUT("/:id", bottleHandler.HandleUpdateBottle)
    bottles.DELETE("/:id", bottleHandler.HandleDeleteBottle)

    // 特殊查询
    bottles.GET("/random", bottleHandler.HandleGetRandomBottles)
    bottles.GET("/hot", bottleHandler.HandleGetHotBottles)
    bottles.GET("/viewed", bottleHandler.HandleGetViewedBottles)
    bottles.GET("/recent-viewed", bottleHandler.HandleGetRecentViewedBottles)
  }

  // TODO: 话题相关路由
  //topics := authenticated.Group("/topics")
  //{
  //  // 待实现...
  //}
}
