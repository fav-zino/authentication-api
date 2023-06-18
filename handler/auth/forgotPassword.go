package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"user_management_system/config"
	"user_management_system/db"
	model "user_management_system/models"
	"user_management_system/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)



func ForgotPasswordHandler(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	
	existingUser := model.User{}
	filter := bson.M{"email":  requestBody.Email }
	err:=db.UserCollection.FindOne(context.Background(),filter).Decode(&existingUser)
	 if err == mongo.ErrNoDocuments {
		  c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"Account not found"})
		  return
	  } else if err != nil{
		  c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message": err})
		  return
	  }



	// Generate a reset token
	resetToken,err := service.GenerateResetToken(existingUser.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	resetLink  := fmt.Sprintf("%s/auth/reset-password/%s",config.AppConfig.ApiURL,resetToken)

	log.Print(config.AppConfig.ApiURL)

	err=service.SendResetPasswordEmail(requestBody.Email, resetLink,existingUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status":"success","message": "A password reset email has been sent to your email"})
}



