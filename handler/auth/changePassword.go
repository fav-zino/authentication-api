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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)




func ChangePasswordHandler(c *gin.Context){

  var requestBody struct{
    TokenString string `json:"token" binding:"required"`
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required"`
  }


  err:= c.BindJSON(&requestBody); if err != nil{
	  c.IndentedJSON(http.StatusBadRequest ,gin.H{"status":"error","message": err})
    return
  }
  if requestBody.OldPassword == ""  {
	c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Missing required field: 'old_password'"})
	return
  }

  if requestBody.TokenString == "" {
	c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Missing required field: 'token'"})
	return
  }

  if requestBody.NewPassword == "" {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "Missing required field: 'new_password'"})
    return
  }

  if len(requestBody.NewPassword) < 5 {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"New Password too short, should be at least 5 characters long"})
    return
  }

  userIDString,err := service.ValidateTokenAndReturnUserID(requestBody.TokenString)
     if err != nil {
		c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message":err})
		return
	}

	//Convert the ID string to a primitive.ObjectID value
    userID, _ := primitive.ObjectIDFromHex(userIDString)

  //get old password hash from mongo
  existingUser := model.User{}
  filter := bson.M{"_id":  userID }
  err=db.UserCollection.FindOne(context.Background(),filter).Decode(&existingUser)
   if err == mongo.ErrNoDocuments {
        c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"Account not found"})
        return
    } else if err != nil{
        c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message": err})
        return
    }

  //Compare the old password with the hash stored in mongo
  err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password),[]byte(requestBody.OldPassword))
  if err !=nil {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"Please enter the old password correctly"})
	return
  }

   //hash new password
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
   if err !=nil {
	 c.IndentedJSON(http.StatusInternalServerError,gin.H{"status":"error","message":err})
	 return
   }

  //update the password field
  update := bson.M{
	"$set": bson.M{"password": string(hashedPassword)},
	}
	filter = bson.M{"_id":  userID }
	var result model.User
	if err = db.UserCollection.FindOneAndUpdate(context.Background(), filter,update).Decode(&result); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"status":"error","message": err})
		return
	}


  c.IndentedJSON(http.StatusOK,gin.H{"status":"success","message":"Password update successful"})

}
