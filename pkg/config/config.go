package config

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/banyar/go-packages/pkg/frontlog"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Configuration struct {
	Logging struct {
		File struct {
			Enabled         bool
			Name            string
			TimestampInName bool
		}
		Level string
	}
}

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

func LoadConfig() *Configuration {

	wd, _ := os.Getwd()
	fmt.Println("Current Working Directory:", wd)
	configPath := filepath.Join(wd, "config")
	fmt.Println("Expected config path:", configPath)

	if _, err := os.Stat(filepath.Join(configPath, "config.json")); os.IsNotExist(err) {
		fmt.Println("config.json does NOT exist in config directory")
	} else {
		fmt.Println("config.json exists in config directory")
	}

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error loading config file: %w", err))
	}

	var cfg Configuration
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unable to decode config into struct: %w", err))
	}

	return &cfg
}

func (c *Configuration) GetLoggingConfig() (frontlog.LogConfig, string) {
	return frontlog.LogConfig{
		Enabled:         viper.GetBool("logging.file.enabled"),
		Name:            viper.GetString("logging.file.name"),
		TimestampInName: viper.GetBool("logging.file.timestamp_in_name"),
	}, viper.GetString("logging.level")
}
