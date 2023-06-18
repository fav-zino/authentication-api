package auth

import (
	"context"
	"net/http"
	"user_management_system/db"
	model "user_management_system/models"
	"user_management_system/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ResetPasswordRenderHtmlHandler(c *gin.Context) {
    resetToken := c.Param("reset-token")
    _, err := service.ValidateTokenAndReturnUserID(resetToken)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset token"})
        return
    }
	data := gin.H{
		"resetToken": resetToken,
	}
	c.HTML(http.StatusOK, "reset_password.html", data)   
}

func ResetPasswordHandler(c *gin.Context) {
	resetToken := c.Param("reset-token")

	userIDString, err := service.ValidateTokenAndReturnUserID(resetToken)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset token"})
        return
    }

	newPassword := c.PostForm("newPassword") 
	confirmPassword := c.PostForm("confirmPassword") 
  
	if newPassword == "" || confirmPassword == "" {
	  c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "password password must be entered"})
	  return
	}
	
  
	if newPassword !=  confirmPassword  {
	  c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Password and confirm password must match"})
	  return
	}

  
	if (len(newPassword) < 5) || (len(confirmPassword) < 5)  {
	  c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"Password too short, should be at least 5 characters long"})
	  return
	}
  
	  //Convert the ID string to a primitive.ObjectID value
	  userID, _ := primitive.ObjectIDFromHex(userIDString)
  
  
	 //hash new password
	 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	 if err !=nil {
	   c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message":err})
	   return
	 }
  
	//update the password field
	update := bson.M{
	  "$set": bson.M{"password": string(hashedPassword)},
	  }
	  filter := bson.M{"_id":  userID }
	  var result model.User
	  if err = db.UserCollection.FindOneAndUpdate(context.Background(), filter,update).Decode(&result); err != nil {
		  c.IndentedJSON(http.StatusInternalServerError, gin.H{"status":"error","message": err})
		  return
	  }

  

    c.IndentedJSON(http.StatusOK, gin.H{"status":"success","message": "Password reset successful"})
}
