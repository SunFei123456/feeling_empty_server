package router

import (
  "github.com/labstack/echo/v4"
  "quick-start/handle"
)

var userHandle = handle.UserHandler{}
var articleHandle = handle.ArticleHandler{}

var progressHandle = handle.ProgressHandler{}

var websiteHandle = handle.WebsiteHandler{}

// 统一管理路由
func SetupRoutes(e *echo.Echo) {
  // 加个api/v1
  apiV1 := e.Group("/api/v1")

  // 在 "api/v1" 路由组中定义路由
  apiV1.POST("/user/create", userHandle.Create)
  apiV1.GET("/user/:id", userHandle.Get)
  apiV1.GET("/user/:id/articles_total", userHandle.GetArticlesTotal)

  apiV1.POST("/article/create", articleHandle.Create)
  apiV1.GET("/article/category/:category", articleHandle.Get)
  apiV1.GET("/article/:id", articleHandle.GetOne)
  apiV1.GET("/article/latest", articleHandle.GetLatest)
  apiV1.GET("/article/category/:category/list", articleHandle.GetArticlesListByCategory)

  apiV1.GET("/progress/list", progressHandle.Get)

  apiV1.GET("/website/:category", websiteHandle.Get)
}
