package jwthandler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenerateJWT(username string) (string, error) {
	var jwtKey = []byte("my_secret_key")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateJWT(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})
	// How we get the user from token!!!
	tclaims, _ := token.Claims.(jwt.MapClaims)
	fmt.Println(tclaims["user_id"])
	//----------------------------------------
	if err != nil {
		return "", err
	}
	if token.Valid {
		return "OK", nil
	}
	return "", nil
}
