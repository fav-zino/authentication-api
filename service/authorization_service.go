package service

import (
	"authentication_api/config"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokenFromID(_id interface{}) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["_id"] = _id
	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(config.AppConfig.TokenSecret))
	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func ValidateTokenAndReturnUserID(requestToken string) (string, error) {
	// Parse the JWT token using the secret key
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the token secret as the key
		return []byte(config.AppConfig.TokenSecret), nil
	})

	if err != nil {
		return "", err
	}
	// Verify the token is valid and has not expired
	if !token.Valid {
		return "", err
	}
	// Get the user ID from the "_id" claim
	userID, ok := token.Claims.(jwt.MapClaims)["_id"].(string)
	if !ok {
		return "", err
	}

	return userID, nil
}
