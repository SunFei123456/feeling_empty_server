package tools

import (
  "fmt"
  _ "github.com/labstack/echo/v4"
  "github.com/labstack/gommon/log"
  "os"
  "reflect"
  "strconv"
)

// 判断是否是生产环境
func IsProduction() bool {
  return os.Getenv("GO_ENV") == "production"
}

func ToMap(entity any, fields ...string) map[string]any {
  resultMap := make(map[string]any)
  val := reflect.ValueOf(entity)
  // If a pointer is passed, get the value that the pointer points to
  if val.Kind() == reflect.Ptr {
    val = val.Elem()
  }

  typ := val.Type()

  for i := 0; i < val.NumField(); i++ {
    field := typ.Field(i)
    fieldVal := val.Field(i)

    // 如果field是嵌入字段，将其提取出来作为一级字段
    if field.Anonymous {
      // 递归调用ToMap函数
      embeddedResult := ToMap(fieldVal.Interface(), fields...)
      for jsonTag, value := range embeddedResult {
        resultMap[jsonTag] = value
      }
    } else {
      jsonTag := field.Tag.Get("json")

      // Check if field is in the list of fields to include
      for _, name := range fields {
        if jsonTag == name {
          resultMap[jsonTag] = val.Field(i).Interface()
          break
        }
      }
    }
  }
  return resultMap
}

// 解析页码 + 边界检查
func ParsePageAndCheckParam(pageParam string) (int, error) {
  // 页码为空, 默认值给1(首页)
  if pageParam == "" {
    return 1, nil
  }
  page, err := strconv.Atoi(pageParam)
  // 如果转换不成功，默认返回1 (首页)
  if err != nil {
    log.Errorf("无效的页码数: %v", pageParam)
    return 1, fmt.Errorf("无效的页码数: %v", pageParam)
  }

  // 处理最小边界值
  if page < 1 {
    log.Errorf("页码数不能为0或负数: %v", pageParam)
    return 1, fmt.Errorf("页码数不能为0或负数: %v", pageParam)
  }

  return page, nil
}
