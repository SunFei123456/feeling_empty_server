package tools

import (
  _ "github.com/labstack/echo/v4"
  "os"
  "reflect"
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
