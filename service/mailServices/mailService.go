package mailServices

import (
	"bytes"
	"codeSandbox/model"
	"codeSandbox/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
)

type MailService struct {
}

var SERVER_HOST = utils.Config.Server.Host
var EMAIL_FROM = utils.Config.Server.Email.From
var EMAIL_PWD = utils.Config.Server.Email.Password
var EMAIL_TO = utils.Config.Server.Email.To
var EMAIL_PORT = utils.Config.Server.Email.Port
var EMAIL_HOST = utils.Config.Server.Email.Host

var MailServiceInstance MailService

func (mailService *MailService) SendToUser(user *model.User) {
	data := map[string]any{
		"creatAt":  user.CreatedAt.Format("2006-01-02 15:04:05"),
		"host":     SERVER_HOST,
		"username": user.Username,
	}

	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("utils/mail/template/toUser.html")
	if err != nil {
		log.Errorf("create template failed, err:%v", err)
		return
	}
	var message bytes.Buffer
	// 写入数据并转为字符串
	err = tmpl.Execute(&message, data)
	if err != nil {
		log.Errorf("Execute failed, err:%v", err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL_FROM) // 发件人

	m.SetHeader("To", user.Email)            // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetHeader("Bcc", EMAIL_TO)             // 暗送，可以多个
	m.SetHeader("Subject", "代码沙箱平台注册审核通过啦！") // 邮件主题

	m.SetBody("text/html", message.String())

	d := gomail.NewDialer(
		EMAIL_HOST,
		EMAIL_PORT,
		EMAIL_FROM,
		EMAIL_PWD,
	)
	// 关闭SSL协议认证
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = d.DialAndSend(m); err != nil {
		log.Errorf("DialAndSend err:%v", err)
	} else {
		log.Infof("Send to %v mail success!", user.Email)
	}
}

func (mailService *MailService) SendToMyself(ip string) {
	data := map[string]any{
		"loginAt": time.Now().Format("2006-01-02 15:04:05"),
		"ip":      ip,
	}

	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("utils/mail/template/toSelf.html")
	if err != nil {
		log.Errorf("create template failed, err:%v", err)
		return
	}
	var message bytes.Buffer
	// 写入数据并转为字符串
	err = tmpl.Execute(&message, data)
	if err != nil {
		log.Errorf("Execute failed, err:%v", err)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL_FROM) // 发件人

	m.SetHeader("To", EMAIL_TO)               // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetHeader("Subject", "警告！代码沙箱平台后台有人登录") // 邮件主题

	m.SetBody("text/html", message.String())

	d := gomail.NewDialer(
		EMAIL_HOST,
		EMAIL_PORT,
		EMAIL_FROM,
		EMAIL_PWD,
	)
	// 关闭SSL协议认证
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err = d.DialAndSend(m); err != nil {
		log.Errorf("DialAndSend err:%v", err)
	} else {
		log.Infof("Send to myself %v mail success!", EMAIL_TO)
	}
}
