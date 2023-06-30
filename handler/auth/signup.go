package auth

import (
	"authentication_api/db"
	model "authentication_api/models"
	"authentication_api/service"
	"context"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// @Summary     
// @Description Create a user account
// @Tags 		auth
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} signupResponseBody
// @Param 		signupRequestBody body signupRequestBody true " "
// @Router 		/auth/signup [post]
func SignupHandler(c *gin.Context) {
	var requestBody signupRequestBody
	

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
	if requestBody.Name == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing required field: 'name'"})
		return
	}

	//check email is valid
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(requestBody.Email) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Please enter a valid email"})
		return
	}

	//check if email is already used
	existingUser := model.User{}
	filter := bson.M{"email": requestBody.Email}
	err = db.UserCollection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		c.IndentedJSON(http.StatusConflict, gin.H{"status": "error", "message": "Email is already in use"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	if len(requestBody.Password) < 5 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Password too short, should be at least 5 characters long"})
		return
	}

	//generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	userDoc := bson.M{
		"name":     requestBody.Name,
		"email":    requestBody.Email,
		"password": string(hashedPassword),
	}

	result, err := db.UserCollection.InsertOne(context.Background(), userDoc)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	//generate token
	token, err := service.GenerateTokenFromID(result.InsertedID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err})
		return
	}

	response := signupResponseBody{
		ID:   result.InsertedID,
		Name:  requestBody.Name,
		Email: requestBody.Email,
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "message": "user created successfully", "token": token, "user": response})

}


type signupRequestBody struct {
	Name     string `json:"name" `     //required
	Email    string `bson:"email"`     //required
	Password string `json:"password" ` //required
}

type signupResponseBody struct {
	ID interface{} `json:"_id"`
	Email    string       `json:"email"`      
	Name string `json:"name"`
}