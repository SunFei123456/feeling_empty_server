package handle

import (
  "encoding/json"
  "fmt"
  "github.com/google/uuid"
  "github.com/labstack/echo/v4"
  "net/http"
  "os"
  "quick-start/db"
  "quick-start/model"
  "quick-start/structs"
  "strconv"
  "time"
)

type CommentHandler struct{}

// 发表评论
func (h *CommentHandler) Create(c echo.Context) error {
  // 获取前端传过来的参数
  var comment model.Comment // nickname avatar body site
  if err := c.Bind(&comment); err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "参数错误")
  }
  // 判断输入的nickname 是否已经存在于数据库当中
  var count int64
  if err := db.DB.Model(model.Comment{}).Where("nickname = ?", comment.Nickname).Count(&count).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  if count > 0 {
    return ErrorResponse(c, http.StatusBadRequest, "昵称已存在,请更换重试")
  }

  // 生成访客 ID
  comment.VisitorId = uuid.New().String()

  // 获取请求者的浏览器信息
  comment.UserAgent = c.Request().UserAgent()

  // 获取 IP 地址的地理位置信息
  res, err := getLocationInfo()
  if err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  comment.City = res.City
  // 设置当前的时间
  comment.CreatedAt = time.Now()
  // 更新的时间
  comment.UpdatedAt = time.Now()

  // 将数据写入数据库
  if err := db.DB.Create(&comment).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  return SuccessResponse(c, "创建评论成功")
}

// 辅助函数：获取地理位置信息
func getLocationInfo() (structs.LocationInfo, error) {
  //https://restapi.amap.com/v3/ip?key=08fb53dd80c6d3b2a8f5f322545d567d
  url := os.Getenv("IP_LOCATION_API_URL") + "?key=" + os.Getenv("GAO_DE_IP_PARSE_KEY") // 5,000次/天
  resp, err := http.Get(url)
  if err != nil {
    return structs.LocationInfo{}, err
  }
  defer resp.Body.Close()

  var LocationInfo structs.LocationInfo
  if err := json.NewDecoder(resp.Body).Decode(&LocationInfo); err != nil {
    return structs.LocationInfo{}, err
  }
  fmt.Println("呵呵", LocationInfo)
  return LocationInfo, nil
}

// 根据评论类型获取评论列表
func (h *CommentHandler) Get(c echo.Context) error {
  var comments []model.Comment
  // 获取评论类型参数
  commentableType := c.Param("commentable_type")
  // 获取分页参数
  page, _ := strconv.Atoi(c.QueryParam("page")) // 页码
  // 如果没有传页码，默认第一页
  if page == 0 {
    page = 1
  }
  size := 7 // 每页展示多少条
  // 计算当前页的起始位置
  offset := (page - 1) * size // 计算偏移量

  if err := db.DB.Where("commentable_type = ?", commentableType).Order("created_at DESC").Limit(size).Offset(offset).Find(&comments).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  // 组装分页器数据
  var total int64
  if err := db.DB.Model(&model.Comment{}).Where("commentable_type = ?", commentableType).Count(&total).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }

  // 计算分页信息
  meta := structs.NewPagination(total, page, size)

  return PagedOkResponse(c, comments, meta)
}

// 删除评论
func (h *CommentHandler) Delete(c echo.Context) error {
  id, err := strconv.Atoi(c.Param("id"))
  if err != nil {
    return ErrorResponse(c, http.StatusBadRequest, "参数错误")
  }
  if err := db.DB.Where("id = ?", id).Delete(&model.Comment{}).Error; err != nil {
    return ErrorResponse(c, http.StatusInternalServerError, err.Error())
  }
  return SuccessResponse(c, "删除成功")
}

//
