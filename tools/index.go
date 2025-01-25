package tools

import (
  "fmt"
  _ "github.com/labstack/echo/v4"
  "github.com/labstack/gommon/log"
  "os"
  "reflect"
  "strconv"

  "crypto/rand"
  "math/big"
)

// IsProduction 判断是否是生产环境
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

// ParsePageAndCheckParam 解析页码 + 边界检查
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

// StringToUint 将字符串转换为无符号整数
func StringToUint(s string) uint {
  i, err := strconv.ParseUint(s, 10, 32)
  if err != nil {
    log.Errorf("Error when converting string to uint: %v", err)
    return 0
  }
  return uint(i)
}

// RandomArrayElement 从数组中随机返回一个元素
func RandomArrayElement(arr []string) string {
  // 生成一个随机的索引
  randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(arr))))
  if err != nil {
    // 如果生成随机数失败，可以返回第一个元素或处理错误
    return arr[0]
  }

  // 返回随机索引对应的元素
  return arr[randomIndex.Int64()]
}
