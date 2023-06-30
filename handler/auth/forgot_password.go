package auth

import (
	"authentication_api/config"
	"authentication_api/db"
	model "authentication_api/models"
	"authentication_api/service"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary 	
// @Description Forgot user password
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Param 		forgotPasswordRequestBody body forgotPasswordRequestBody true " "
// @Router /auth/forgot-password [post]
func ForgotPasswordHandler(c *gin.Context) {
	var requestBody forgotPasswordRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	existingUser := model.User{}
	filter := bson.M{"email": requestBody.Email}
	err := db.UserCollection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Account not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	// Generate a reset token
	resetToken, err := service.GenerateResetToken(existingUser.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	resetLink := fmt.Sprintf("%s/auth/reset-password/%s", config.AppConfig.ApiURL, resetToken)

	log.Print(resetLink)

	err = service.SendResetPasswordEmail(requestBody.Email, resetLink, existingUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "A password reset email has been sent to your email"})
}

type forgotPasswordRequestBody struct {
	Email    string `json:"email"`
}


