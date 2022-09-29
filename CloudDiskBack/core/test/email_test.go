package test

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Cloud Disk <onemorelight234@163.com>"
	e.To = []string{"501124524@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>123456</h1>")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "onemorelight234@163.com", "RBQKKMWRKVAUISFS", "smtp.163.com"),
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "smtp.163.com",
		})
	if err != nil {
		t.Fatal(err)
	}

}
