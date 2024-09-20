package structs

type LocationInfo struct {
  Status    string `json:"status"`    // 状态码
  Info      string `json:"info"`      // 状态信息
  Infocode  string `json:"infocode"`  // 信息码
  Province  string `json:"province"`  // 省份名称
  City      string `json:"city"`      // 城市名称
  Adcode    string `json:"adcode"`    // 区域编码
  Rectangle string `json:"rectangle"` // 矩形范围坐标（经纬度）
}

type Pagination struct {
  Count int64 `json:"count"`
  Page  int   `json:"page"`
  Items int   `json:"items"`
  Prev  int   `json:"prev"`
  Next  int   `json:"next"`
  Last  int   `json:"last"`
}

func NewPagination(total int64, page int, limit int) Pagination {
  last := (int(total) + limit - 1) / limit // 计算最后一页
  prev := page - 1                         // 计算上一页
  next := page + 1                         // 计算下一页
  if prev < 1 {
    prev = 1
  } else if prev > last {
    prev = last
  }
  if next > last {
    next = last
  }
  return Pagination{
    Count: total,
    Page:  page,
    Items: limit,
    Prev:  prev,
    Next:  next,
    Last:  last,
  }
}
