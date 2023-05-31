package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gookit/config/v2"
	jsonDriver "github.com/gookit/config/v2/json"
	"log"
	"time"
)

func genApiKey() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "circle-connector",
		Subject:   "circle-stats-aggregator",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        "3",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(config.String("keySecret")))

	fmt.Println(tokenString)
}

func main() {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(jsonDriver.Driver)

	err := config.LoadFiles("config.json")
	if err != nil {
		log.Panic(err)
	}

	genApiKey()
}
