package auth

import (
	"context"
	"net/http"
	"regexp"
	"user_management_system/db"
	model "user_management_system/models"
	"user_management_system/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)




func SignupHandler(c *gin.Context){
  var requestBody struct {
    Name     string `json:"name" binding:"required"`
    Email    string `bson:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=5"`
  }

  err:= c.BindJSON(&requestBody); if err != nil{
	  c.IndentedJSON(http.StatusBadRequest ,gin.H{"status":"error","message": err})
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
	if  requestBody.Name == "" {
	  c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Missing required field: 'name'"})
	  return
	}

  //check email is valid
  emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
  if !emailRegex.MatchString(requestBody.Email) {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Please enter a valid email"})
    return
  }

  //check if email is already used
  existingUser := model.User{}
  filter := bson.M{"email":  requestBody.Email }
  err=db.UserCollection.FindOne(context.Background(),filter).Decode(&existingUser)
   if err == nil {
        c.IndentedJSON(http.StatusConflict,gin.H{"status":"error","message":"Email is already in use"})
        return
    } else if err != mongo.ErrNoDocuments {
        c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message": err})
        return
    }

  if len(requestBody.Password) < 5 {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"Password too short, should be at least 5 characters long"})
    return
  }


  //generate hashed password
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
  if err !=nil {
    c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message":err})
    return
  }

  userDoc := bson.M{
    "name":requestBody.Name,
    "email":requestBody.Email,
    "password":string(hashedPassword),
  }


  result,err:=db.UserCollection.InsertOne(context.Background(),userDoc); if err != nil{
    c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message":err})
    return
  }

  //generate token
  token,err := service.GenerateTokenFromID(result.InsertedID)
     if err != nil {
        c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message":err})
        return
    }



  response := map[string]interface{}{
    "_id":result.InsertedID,
    "name":requestBody.Name,
    "email":requestBody.Email,
  }

  c.IndentedJSON(http.StatusOK,gin.H{"status":"success","message":"user created successfully","token":token,"user":response})

}



