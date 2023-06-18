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

type ChangePassword struct{
	TokenString string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}


func ChangePasswordHandler(c *gin.Context){
  var changePasswordBody ChangePassword
  err:= c.BindJSON(&changePasswordBody); if err != nil{
	  c.IndentedJSON(http.StatusBadRequest ,gin.H{"status":"error","message": "some error occured"})
    return
  }
  if changePasswordBody.OldPassword == "" || changePasswordBody.TokenString == "" {
	c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "all required fields must be filled"})
	return
  }

  if changePasswordBody.NewPassword == "" {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"message": "New password must be entered"})
    return
  }

  if len(changePasswordBody.NewPassword) < 5 {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"New Password too short, should be at least 5 characters long"})
    return
  }

  userIDString,err := service.ValidateTokenAndReturnUserID(changePasswordBody.TokenString)
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
  err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password),[]byte(changePasswordBody.OldPassword))
  if err !=nil {
    c.IndentedJSON(http.StatusBadRequest,gin.H{"status":"error","message":"Please enter the old password correctly"})
	return
  }

   //hash new password
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordBody.NewPassword), bcrypt.DefaultCost)
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
