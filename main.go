package main

import (
	// "html/template"
	"log"
	"os"
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
	startGin()
}


func startGin()  {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.LoadHTMLGlob("view/*.html")
	//  templateEngine := template.New("")
	//  templateEngine.ParseGlob("view/*.html")
	//  router.SetHTMLTemplate(templateEngine)


	routes.LoadAuthRoutes(router)

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
        log.Panicf("Error starting server: %s", err)
	}
	// gin.SetMode(gin.DebugMode)
	// err := router.Run("localhost:8080")
	// if err != nil {
	// 	log.Fatal("Error starting server:", err)
	// }
}
