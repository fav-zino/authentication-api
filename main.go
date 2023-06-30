package main

import (
	"authentication_api/config"
	"authentication_api/db"
	_ "authentication_api/docs"
	"authentication_api/routes"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

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


//documentation at /docs/index.html


// @title Authentication API
// @desciption Authentication API using Gin
// @host localhost:3000
func main() {
	startGin()
}

func startGin() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.LoadHTMLGlob("view/*.html")
	//  templateEngine := template.New("")
	//  templateEngine.ParseGlob("view/*.html")
	//  router.SetHTMLTemplate(templateEngine)

	// Serve Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
