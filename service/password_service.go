package service

import (
	// "fmt"
	"log"
	"strconv"
	"time"
	"user_management_system/config"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gomail.v2"
)



func SendResetPasswordEmail(toEmail string, resetLink string,userName string) error {

    messageBody := 
    `<h1>Reset Password</h1>
    <p>Hi `+userName+ `,</p>
    <p>We received a request to reset your password. To proceed, please click the link below:</p>
    <p><a href="`+resetLink+ `" style="color:blue;">Reset Password</a></p>
    <p>If you didn't request a password reset, please ignore this email.</p>
    <p>Best regards,</p>
    <p>Kuro team</p>`

    // Initialize the email message
    message := gomail.NewMessage()
    message.SetHeader("From", config.AppConfig.EmailUsername)
    message.SetHeader("To", "ewomazinofavour22@gmail.com")
    message.SetHeader("Subject", "Reset Password")
    message.SetBody("text/html",messageBody)
    
    // Create a new email sending session
	
    port, _ := strconv.Atoi(config.AppConfig.EmailSMTPPort)
    
    d := gomail.NewDialer(config.AppConfig.EmailSMTPHost, port, config.AppConfig.EmailUsername, config.AppConfig.EmailPassword)
    
    
    // Send the email
    if err := d.DialAndSend(message); err != nil {
        log.Print(err)
        return err
    }

    return nil
}


func GenerateResetToken(userID string) (string, error) {
    // Create a new JWT token
    token := jwt.New(jwt.SigningMethodHS256)

    // Set the claims
    claims := token.Claims.(jwt.MapClaims)
    claims["_id"] = userID
    claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // Expiration time is set to 15 minutes from now

    // Sign the token with a secret key
    tokenString, err := token.SignedString([]byte(config.AppConfig.TokenSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}





