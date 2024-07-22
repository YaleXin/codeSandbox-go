package mailServices

import (
	"bytes"
	"codeSandbox/model"
	"codeSandbox/utils"
	"codeSandbox/utils/global"
	"codeSandbox/utils/tool"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"html/template"
	"strconv"
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
var REGISTER_URL_KEY = utils.Config.Server.RegisterUrlKey
var MailServiceInstance MailService

func (mailService *MailService) SendToUser(user *model.User) int {
	data := map[string]any{
		"creatAt":     user.CreatedAt.Format("2006-01-02 15:04:05"),
		"host":        SERVER_HOST,
		"username":    user.Username,
		"registerUrl": generateRegisterUrl(user.ID),
	}

	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("utils/mail/template/toUser.html")
	if err != nil {
		log.Errorf("create template failed, err:%v", err)
		return global.SYSTEM_ERROR
	}
	var message bytes.Buffer
	// 写入数据并转为字符串
	err = tmpl.Execute(&message, data)
	if err != nil {
		log.Errorf("Execute failed, err:%v", err)
		return global.SYSTEM_ERROR
	}
	subject := "代码沙箱平台注册审核通过啦！"
	code := sendEmail(user.Email, message.String(), subject, true)
	return code
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
	subject := "警告！代码沙箱平台后台有人登录"
	sendEmail(EMAIL_TO, message.String(), subject, false)
}
func (mailService *MailService) SendToMyselfToAuditUser(user *model.User) {
	data := map[string]any{
		"username": user.Username,
		"email":    user.Email,
	}

	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("utils/mail/template/toSelfToAuditUser.html")
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
	subject := "代码沙箱平台有新用户注册啦"
	sendEmail(EMAIL_TO, message.String(), subject, false)
}
func generateRegisterUrl(userId uint) string {
	userIdStr := strconv.Itoa(int(userId))
	// 注册链接为 http://code.yalexin.top/api/v1/user/check?token=KEY
	// KEY = <base64(userid)>.<md5(key+userid)>
	encodedBytes := base64.StdEncoding.EncodeToString([]byte(userIdStr))
	md5Str := tool.MD5Str(REGISTER_URL_KEY + userIdStr)
	return global.REGISTER_URL_PREFIX + encodedBytes + "." + md5Str
}

// sendEmail
//
//	@Description:
//	@param toEmail 收件人地址
//	@param htmlStr 邮件 html 内容
//	@param subject 邮件主题
//	@param bccEnable 是否抄送
func sendEmail(toEmail, htmlStr, subject string, bccEnable bool) int {

	m := gomail.NewMessage()
	m.SetHeader("From", EMAIL_FROM) // 发件人

	m.SetHeader("To", toEmail)      // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接
	m.SetHeader("Subject", subject) // 邮件主题
	if bccEnable {
		m.SetHeader("Bcc", EMAIL_TO)
	}

	m.SetBody("text/html", htmlStr)

	d := gomail.NewDialer(
		EMAIL_HOST,
		EMAIL_PORT,
		EMAIL_FROM,
		EMAIL_PWD,
	)
	// 关闭SSL协议认证
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Errorf("DialAndSend err:%v", err)
		return global.SYSTEM_ERROR
	} else {
		log.Infof("Send to  %v mail success!", EMAIL_TO)
		return global.SUCCESS
	}
}
