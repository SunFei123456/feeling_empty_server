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
  bottleViewHandler := handler.NewBottleViewHandler(db.DB)
  cosHandler := handler.COSHandler{}
  bottleInteractionHandler := handler.NewBottleInteractionHandler(db.DB)
  oceanHandler := handler.NewOceanHandler(db.DB)
  topicHandler := handler.NewTopicHandler(db.DB)

  // 关注
  var userFolloweesHandler = handler.UserFolloweesHandler{}
  // 粉丝
  var userFollowersHandler = handler.UserFollowersHandler{}

  // API 路由组
  api := e.Group("/api/v1")

  // 公开路由组
  auth := api.Group("/auth")
  {
    auth.POST("/register", userHandler.HandleRegister)
    auth.POST("/login", userHandler.HandleLogin)
    auth.POST("/qq-email-login", userHandler.HandleQQEmailLogin)
    auth.POST("/send-code", userHandler.HandleSendEmailCode)
  }

  // 需要认证的路由组
  authenticated := api.Group("")
  authenticated.Use(middleware.JWT())

  // 用户相关路由
  users := authenticated.Group("/users")
  {
    // 用户信息
    users.GET("", userHandler.HandleGetCurrentUser)
    // 更新用户信息
    users.PUT("", userHandler.HandleUpdateCurrentUser)

    // 根据user_id 获取不同的用户信息
    users.GET("/:user_id", userHandler.HandleGetUserByID)
  }

  // 漂流瓶相关路由
  bottles := authenticated.Group("/bottles")
  {
    // 基础操作
    bottles.POST("", bottleHandler.HandleCreateBottle)
    bottles.GET("", bottleHandler.HandleGetBottles)
    bottles.PUT("/:id", bottleHandler.HandleUpdateBottle)
    bottles.DELETE("/:id", bottleHandler.HandleDeleteBottle)

    // 特殊查询
    bottles.GET("/random", bottleHandler.HandleGetRandomBottles)
    bottles.GET("/hot", bottleHandler.HandleGetHotBottles)
    bottles.GET("/recent-viewed", bottleHandler.HandleGetRecentViewedBottles)

    // 互动相关路由
    bottles.POST("/:id/resonate", bottleInteractionHandler.HandleResonateBottle)
    bottles.DELETE("/:id/resonate", bottleInteractionHandler.HandleCancelResonateBottle)
    bottles.GET("/resonated", bottleInteractionHandler.HandleGetUserResonatedBottles)

    bottles.POST("/:id/favorite", bottleInteractionHandler.HandleFavoriteBottle)
    bottles.DELETE("/:id/favorite", bottleInteractionHandler.HandleCancelFavoriteBottle)
    bottles.GET("/favorited", bottleInteractionHandler.HandleGetUserFavoriteBottles)

    bottles.POST("/:id/share", bottleInteractionHandler.HandleShareBottle)

    bottles.GET("/:id/interaction", bottleInteractionHandler.HandleGetBottleInteractionStatus)
  }

  // 漂流瓶浏览记录相关路由
  bottleViews := authenticated.Group("/bottle-views")
  {
    // 获取漂流瓶浏览记录
    bottleViews.GET("", bottleViewHandler.HandleGetBottleViews)
    // 删除指定漂流瓶浏览记录
    bottleViews.DELETE("/:id", bottleViewHandler.HandleDeleteBottleView)
    // 删除用户的全部漂流瓶浏览记录
    bottleViews.DELETE("", bottleViewHandler.HandleDeleteAllBottleViews)
    // 创建漂流瓶浏览记录
    bottleViews.POST("", bottleViewHandler.HandleCreateBottleView)
  }

  // TencentCOS 相关路由
  tcos := authenticated.Group("/cos")
  {
    tcos.GET("/upload-token", cosHandler.GetUploadToken)
  }

  // 海域相关路由
  oceans := authenticated.Group("/oceans")
  {
    oceans.GET("", oceanHandler.HandleGetOceans)                         // 获取所有海域信息
    oceans.GET("/:ocean_id/bottles", oceanHandler.HandleGetOceanBottles) // 获取指定海域的瓶子
  }

  // 话题相关路由
  topics := authenticated.Group("/topics")
  {
    topics.GET("/system", topicHandler.HandleGetSystemTopics)      // 获取系统话题
    topics.GET("/:id/bottles", topicHandler.HandleGetTopicBottles) // 获取话题下的漂流瓶
    topics.GET("/:id", topicHandler.HandleGetTopicInfo)            // 获取话题详情
    topics.GET("/hot", topicHandler.HandleGetHotTopics)            // 获取热门话题
    topics.POST("", topicHandler.HandleCreateTopic)                // 创建话题
    topics.GET("", topicHandler.HandleGetAllTopics)                // 获取所有话题
    topics.GET("/search", topicHandler.HandleSearchTopics)         // 搜索话题
  }

  // 用户关注模块
  follows := authenticated.Group("/followees")
  {
    follows.GET("/user/:id", userFolloweesHandler.Index)                         //  获取指定用户的关注列表
    follows.GET("/user/:id/follow_status", userFolloweesHandler.GetFollowStatus) //  获取两个用户之间的状态()
    follows.POST("/user/:id/follow", userFolloweesHandler.FollowUser)            //  关注
    follows.POST("/user/:id/unfollow", userFolloweesHandler.UnfollowUser)        // 取关
  }
  // 用户粉丝模块
  fans := authenticated.Group("/followers")
  {
    fans.GET("/user/:id", userFollowersHandler.Index)                             //  获取指定用户的粉丝列表
    fans.GET("/user/:id/recent", userFollowersHandler.GetRecentThreeDaysFansList) //  获取近三天内新增的粉丝 的列表
  }

}
