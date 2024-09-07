package model

// 定义一个Website 结构体,用来表示website表
type Website struct {
  BaseModel
  Name     string `json:"name"`
  Desc     string `json:"desc"`
  Href     string `json:"href"`
  Logo     string `json:"logo"`
  Tags     string `json:"tags"`
  Category string `json:"category"`
}
