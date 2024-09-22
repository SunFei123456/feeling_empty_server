package router

import (
  "github.com/labstack/echo/v4"
  "quick-start/handle"
)

var userHandle = handle.UserHandler{}
var articleHandle = handle.ArticleHandler{}

var progressHandle = handle.ProgressHandler{}

var websiteHandle = handle.WebsiteHandler{}

var poemHandle = handle.PoemHandler{}

var commentHandle = handle.CommentHandler{}

var likeHandler = handle.LikeHandler{}

var statsHandler = handle.StatsHandler{}

var imageHandler = handle.ImageHandler{}

var cognitiveExpansionHandler = handle.CognitiveExpansionHandler{}

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
  apiV1.GET("/article/hot", articleHandle.GetRandomHotArticles)
  apiV1.PUT("/article/:id/views", articleHandle.IncrementViewCount)
  apiV1.GET("/article/:id/comments", articleHandle.GetCommentsByArticleID)

  apiV1.GET("/progress/list", progressHandle.Get)

  apiV1.GET("/website/:category", websiteHandle.Get)
  apiV1.GET("/website/:category/tags", websiteHandle.GetTagsByCategory)
  apiV1.GET("/website/tag/search", websiteHandle.SearchByTag)

  apiV1.GET("/poem/list", poemHandle.Get)

  // 评论
  apiV1.POST("/comment/create", commentHandle.Create)
  apiV1.GET("/comment/:commentable_type/list", commentHandle.Get)
  apiV1.DELETE("/comment/:id", commentHandle.Delete)

  // 切换喜欢
  apiV1.POST("/reaction/:type/:id/:reaction_type", likeHandler.ToggleLike)

  // 数据统计
  apiV1.GET("/stats/article", statsHandler.GetArticleStats)

  // 上传文件
  apiV1.POST("/upload", imageHandler.Upload)
  // 获取bing 每日壁纸
  apiV1.GET("/bing_wallpaper", imageHandler.GetRandomBingImage)

  // 认知扩展模块

  // 根据id获取指定文章
  apiV1.GET("/cognitive_expansion/article/:id", cognitiveExpansionHandler.GetArticleById)
  // 根据tag获取指定文章列表
  apiV1.GET("/cognitive_expansion/article/tag", cognitiveExpansionHandler.GetArticlesByTag)
  // 获取所有文章tags
  apiV1.GET("/cognitive_expansion/article/tags", cognitiveExpansionHandler.GetTags)
  // 获取最近的10篇文章
  apiV1.GET("/cognitive_expansion/article/latest", cognitiveExpansionHandler.GetLatestArticles)
  // 根据title 全文搜素 模糊匹配
  apiV1.GET("/cognitive_expansion/article/search", cognitiveExpansionHandler.SearchByTitle)
}
