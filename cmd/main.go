package main

import (
	"fmt"

	"github.com/melnikdev/book-mail/config"
	"github.com/melnikdev/book-mail/internal/services/mail"
)

func main() {

	config := config.MustLoad()
	mail := mail.New(config)

	if err := mail.Send(); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}

}
