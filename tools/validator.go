package tools

import (
  "github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
  validator *validator.Validate
}

// NewCustomValidator 创建自定义验证器
func NewCustomValidator() *CustomValidator {
  validate := validator.New()

  // 添加自定义验证规则
  validate.RegisterValidation("custom_rule", func(fl validator.FieldLevel) bool {
    // 自定义验证逻辑
    return true
  })

  return &CustomValidator{validator: validate}
}

// Validate 实现 echo.Validator 接口
func (cv *CustomValidator) Validate(i interface{}) error {
  return cv.validator.Struct(i)
}
