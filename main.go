package main

import (
	"html/template"
	"log"
	"user_management_system/config"
	"user_management_system/db"
	"user_management_system/routes"
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


	routes.LoadAuthRoutes(router)


	

	gin.SetMode(gin.DebugMode)
	err := router.Run("localhost:8080")

	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
