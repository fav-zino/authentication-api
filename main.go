package main

import (
	"html/template"
	"log"
	"user_management_system/config"
	"user_management_system/db"
	"user_management_system/handler/auth"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	config.AppConfig.LoadFromEnv()

	dbErr := db.ConnectToDB()
	if dbErr != nil {
		log.Fatal("Error connecting to database:", dbErr)
	}
}

func main() {
	router := gin.Default()

	 templateEngine := template.New("")
	 templateEngine.ParseGlob("view/*.html")
	 router.SetHTMLTemplate(templateEngine)

	  // Define a route group for allow route to have a common prefix
	//   router := r.Group("/api")

	router.POST("/signup", auth.SignupHandler)
	router.POST("/login", auth.LoginHandler)
	router.POST("/change-password", auth.ChangePasswordHandler)
	router.POST("/forgot-password", auth.ForgotPasswordHandler)
	//render html to reset password
	router.GET("/reset-password/:reset-token", auth.ResetPasswordRenderHtmlHandler)
	//reset the password
	router.POST("/reset-password/:reset-token", auth.ResetPasswordHandler)


	

	gin.SetMode(gin.DebugMode)
	err := router.Run("localhost:8080")

	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
