package account_link

import (
	"connector-backend/models"
	"github.com/gofiber/fiber/v2"
	"log"
)

func SetupRoutes(router fiber.Router) {
	router.Post("/link", LinkAccount)
	router.Get("/steam/all", GetAllAccountsBySteamId)
	router.Get("/steam/:steamId", GetAccountBySteamId)
	router.Get("/discord/all", GetAllAccountsByDiscordId)
	router.Get("/discord/:discordId", GetAccountByDiscordId)

	if err := models.Db.AutoMigrate(&Account{}); err != nil {
		log.Panic(err)
	}
}
