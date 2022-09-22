package main

import (
	account_link "connector-backend/account-link"
	"connector-backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/keyauth/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gookit/config/v2"
	jsonDriver "github.com/gookit/config/v2/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	errMissing = &models.Body{ErrorField: models.Error{
		Status:  401,
		Message: "Missing API Key",
	}}
	errInvalid = &models.Body{ErrorField: models.Error{
		Status:  401,
		Message: "Invalid API Key",
	}}
)

func main() {
	var err error
	models.Db, err = gorm.Open(sqlite.Open("connector.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	config.WithOptions(config.ParseEnv)
	config.AddDriver(jsonDriver.Driver)

	err = config.LoadFiles("config.json")
	if err != nil {
		log.Panic(err)
	}

	//genApiKey()

	app := fiber.New()
	app.Use(cors.New(
		cors.Config{
			Next:             nil,
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
			AllowHeaders:     "",
			AllowCredentials: true,
			ExposeHeaders:    "",
			MaxAge:           0,
		},
	))
	app.Use(logger.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(keyauth.New(keyauth.Config{
		ErrorHandler: errHandler,
		Validator:    validator,
		ContextKey:   "apiKey",
	}))

	account_link.SetupRoutes(app.Group("/account-link/"))

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}

func errHandler(c *fiber.Ctx, err error) error {
	c.Status(fiber.StatusUnauthorized)

	if err == errMissing {
		return c.JSON(errMissing)
	}

	return c.JSON(errInvalid)
}

func validator(c *fiber.Ctx, s string) (bool, error) {
	if s == "" {
		return false, errMissing
	}

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalid
		}

		return []byte(config.String("keySecret")), nil
	})

	fmt.Println(token, err)

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("claims", claims)

		return true, nil
	}

	return false, errInvalid
}

func genApiKey() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "circle-connector",
		Subject:   "circle-stats-bot",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        "2",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(config.String("keySecret")))

	fmt.Println(tokenString)
}
