package service

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"spm/model"
	"strings"
)

const (
	contentTypeTextPlain = "Content-Type: text/plain; charset=UTF-8"
	contentTypeTextHTML  = "Content-Type: text/html; charset=UTF-8"
)

var emailHost, emailPort, emailUsername, emailPassword string
var emailAuth smtp.Auth

type EmailData struct {
	FirstName string
	Date      string
	Price     string
	Comment   string
	Action    string
}

func init() {
	register(initializeEmailService)
}

func initializeEmailService() {
	emailHost = cfg["email_host"]
	emailPort = cfg["email_port"]
	emailUsername = cfg["email_username"]
	emailPassword = cfg["email_password"]
	emailAuth = smtp.PlainAuth("", emailUsername, emailPassword, emailHost)
	log.Println("Email service initialized.")
}

// general function to send email
func SendMail(sender string, from string, to []string, subject string, contentType string, body string) {
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + sender +
		"<" + from + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail(emailHost+":"+emailPort, emailAuth, from, to, msg)
	if err != nil {
		log.Printf("send mail error: %v\n", err)
	}
}

func EmailConfirm(id uint, appointment model.Appointment, action string) {
	sender := "Gabriel and David"
	from := "admin@gmail.com"
	subject := "Appointment Confirmation"

	// load client info
	user := model.User{}
	user.Retrieve(id)

	emailData := EmailData{}
	emailData.Date = appointment.Time.Format("2006-01-02 15:04")
	emailData.FirstName = user.FirstName
	emailData.Action = action
	emailData.Comment = appointment.Comment
	emailData.Price = model.GetPrice(appointment.Option)
	email := user.Email

	var contentType string
	var content string

	// get template
	tmpl, err := template.ParseFiles("template/confirmation.html")
	if err != nil {
		log.Println("Email error, ", err)
		return
	} else {
		var buf bytes.Buffer

		err = tmpl.Execute(&buf, emailData)
		if err != nil {
			log.Println("Email error, ", err)
			return
		} else {
			contentType = contentTypeTextHTML
			content = buf.String()
		}
	}

	SendMail(sender, from, []string{email}, subject, contentType, content)
}

