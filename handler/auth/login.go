package auth

import (
	"authentication_api/db"
	_ "authentication_api/docs"
	model "authentication_api/models"
	"authentication_api/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// @Summary 	
// @Description Login into user account
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} loginResponseBody
// @Param 		loginRequestBody body loginRequestBody true " "
// @Router 		/auth/login [post]
func LoginHandler(c *gin.Context) {
	var requestBody loginRequestBody

	err := c.BindJSON(&requestBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}
	if requestBody.Email == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing required field: 'email'"})
		return
	}

	if requestBody.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing required field: 'password'"})
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

	response := loginResponseBody{
		ID:   existingUser.ID,
		Name:  existingUser.Name,
		Email: existingUser.Email,
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "message": "login successfully", "token": tokenString, "user": response})

}
type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty" `
}

type loginResponseBody struct {
	ID primitive.ObjectID `json:"_id"`
	Email    string       `json:"email"`      
	Name string `json:"name"`
}
