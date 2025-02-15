package mail

import (
	"fmt"
	"log"
	"strconv"
	"time"

	gomail "gopkg.in/mail.v2"

	"github.com/melnikdev/book-mail/config"
)

type Mail struct {
	Config *config.Config
}

type User struct {
	UserId int `json:"user_id"`
}

func New(config *config.Config) *Mail {
	return &Mail{Config: config}
}

func (m *Mail) sendEmail(u *User) error {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", "recipient1@email.com")
	message.SetHeader("Subject", "Hello from the Mailtrap team")

	// Set email body
	message.SetBody("text/plain", "This is the Test Body id: "+strconv.Itoa(u.UserId))

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(m.Config.Sandbox.Host, m.Config.Sandbox.Port, m.Config.Sandbox.Username, m.Config.Sandbox.Password)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (m *Mail) ListenEvent(ch chan User) {
	for {
		select {
		case user, ok := <-ch:
			if !ok {
				fmt.Println("ch1 is closed. Exiting loop.")
				break
			}
			fmt.Println("Received from ch1:", user)
			if err := m.sendEmail(&user); err != nil {
				log.Fatal("Failed mail send:", err)
			}
		case <-time.After(1 * time.Second):
			fmt.Println("Waiting...")
		}
	}
}
