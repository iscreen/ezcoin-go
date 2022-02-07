package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/jordan-wright/email"
)

func Email(To, subject string, body []byte) error {
	to := strings.Split(To, ",")
	return send(to, subject, body)
}

//@author: [maplepie](https://github.com/maplepie)
//@function: send
//@description: Email发送方法
//@param: subject string, body string
//@return: error
func send(to []string, subject string, body []byte) error {
	from := global.GVA_CONFIG.Email.From
	nickname := global.GVA_CONFIG.Email.Nickname
	secret := global.GVA_CONFIG.Email.Secret
	host := global.GVA_CONFIG.Email.Host
	port := global.GVA_CONFIG.Email.Port
	isSSL := global.GVA_CONFIG.Email.IsSSL
	bcc := strings.Split(global.GVA_CONFIG.Email.Bcc, ",")

	var auth smtp.Auth
	if secret != "" {
		auth = smtp.PlainAuth("", from, secret, host)
	}
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = body
	e.Bcc = bcc
	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
