package notice

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

type EmailNotifier struct{}

func (e *EmailNotifier) Notice(contactInfo string, context string) {

	fmt.Printf("Sending to Email %s: %s\n", contactInfo, context)
	// QQ 邮箱配置
	username := os.Getenv("SMTP_Username")
	password := os.Getenv("SMTP_Password")

	//接收者邮箱列表
	mailTo := []string{contactInfo}
	m := gomail.NewMessage()
	m.SetHeader("From", username)   //发送者腾讯邮箱账号
	m.SetHeader("To", mailTo...)    //接收者邮箱列表
	m.SetHeader("Subject", "通知")    //邮件标题
	m.SetBody("text/html", context) //邮件内容,可以是html
	d := gomail.NewDialer("smtp.qq.com", 465, username, password)
	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println("Email send failure!")
		return
	}
	fmt.Println("Email send successfully!")
}
