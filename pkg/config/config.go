package config

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AppName string

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("config.init()")
	AppName = os.Getenv("APP_NAME")
}

func GetMailTemplete() *template.Template {
	tmpl, err := template.ParseFiles("../template/mail_template.html")
	if err != nil {
		fmt.Printf("error parsing template: %v", err)
	}
	return tmpl
}
