package routes

import (
	"user_management_system/handler/auth"

	"github.com/gin-gonic/gin"
)

func LoadAuthRoutes(router *gin.Engine) {


	authRouter := router.Group("/auth")
	
	authRouter.POST("/signup", auth.SignupHandler)
	authRouter.POST("/login", auth.LoginHandler)
	authRouter.POST("/change-password", auth.ChangePasswordHandler)
	authRouter.POST("/forgot-password", auth.ForgotPasswordHandler)
	//render html to reset password
	authRouter.GET("/reset-password/:reset-token", auth.ResetPasswordRenderHtmlHandler)
	//reset the password
	authRouter.POST("/reset-password/:reset-token", auth.ResetPasswordHandler)

}
