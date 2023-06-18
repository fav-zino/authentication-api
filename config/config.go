package config

import (
    "os"
)

type Config struct {
    DatabaseURL      string
    ApiURL      string
    EmailSMTPHost    string
    EmailSMTPPort    string
    EmailUsername string
	TokenSecret string
    EmailPassword string
}

func (c *Config) LoadFromEnv() {
    c.DatabaseURL =     os.Getenv("DB_URL")
    c.EmailSMTPHost=    os.Getenv("EMAIL_SMTP_HOST")
    c.EmailSMTPPort=    os.Getenv("EMAIL_SMTP_PORT")
    c.TokenSecret= 	  os.Getenv("TOKEN_SECRET")
    c.EmailUsername= os.Getenv("EMAIL_SMTP_USERNAME")
    c.EmailPassword= os.Getenv("EMAIL_SMTP_PASSWORD")
    c.ApiURL=os.Getenv("API_URL")
}

var AppConfig = Config{}


