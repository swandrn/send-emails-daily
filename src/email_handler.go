package main

import (
	"fmt"
	"net/smtp"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
)

var wg sync.WaitGroup

func sendEmail(address string, name string) error {
	defer wg.Done()
	err := godotenv.Load()
	if err != nil {
		return err
	}

	emailToSend := email.NewEmail()
	emailToSend.From = os.Getenv("EMAIL_USER")
	emailToSend.To = []string{address}
	emailToSend.Subject = "Vous avez des tâches à compléter !"
	emailToSend.Text = []byte(
		fmt.Sprintf(`Bonjour %s !
		Vous avez des tâches non complétés !`, name),
	)

	//Change password
	err = emailToSend.Send("smtp.mail.yahoo.com:587", smtp.PlainAuth("", os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_APP_PWD"), "smtp.mail.yahoo.com"))
	if err != nil {
		return err
	}
	return nil
}
