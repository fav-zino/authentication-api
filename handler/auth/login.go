package auth

import (
	"context"
	"net/http"
	"user_management_system/db"
	model "user_management_system/models"
	"user_management_system/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)



func LoginHandler(c *gin.Context) {
  var requestBody struct {
    Email    string `bson:"email" binding:"required,email"`
    Password string `json:"password,omitempty" binding:"required,min=5"`
  }

	err := c.BindJSON(&requestBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}
	if requestBody.Email == "" {
	  c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Missing required field: 'email'"})
	  return
	}

	if  requestBody.Password == "" {
	  c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Missing required field: 'password'"})
	  return
	}

	//check if account exist
	existingUser := model.User{}
	filter := bson.M{"email": requestBody.Email}
	err = db.UserCollection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Account not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	//Compare the entered password with the hash stored in mongo
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(requestBody.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "incorrect email or password"})
		return
	}

	//generate token
	tokenString, err := service.GenerateTokenFromID(existingUser.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	response := map[string]interface{}{
		"_id":   existingUser.ID,
		"name":  existingUser.Name,
		"email": existingUser.Email,
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "message": "login successfully", "token": tokenString, "user": response})

}
