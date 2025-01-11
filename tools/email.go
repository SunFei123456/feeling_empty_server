package tools

import (
  "crypto/tls"
  "fmt"
  "gopkg.in/gomail.v2"
  "os"
  "strconv"
)

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
  // 添加调试信息
  port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
  if err != nil {
    return fmt.Errorf("端口配置错误: %v", err)
  }

  m := gomail.NewMessage()
  m.SetHeader("From", os.Getenv("SMTP_USER"))
  m.SetHeader("To", to)
  m.SetHeader("Subject", subject)
  m.SetBody("text/plain", body)

  d := gomail.NewDialer(
    os.Getenv("MAIL_HOST"),
    port,
    os.Getenv("SMTP_USER"),
    os.Getenv("SMTP_PASSWORD"),
  )

  // SSL配置
  d.TLSConfig = &tls.Config{
    InsecureSkipVerify: true,
    ServerName:         "smtp.qq.com",
  }
  d.SSL = true

  return d.DialAndSend(m)
}
